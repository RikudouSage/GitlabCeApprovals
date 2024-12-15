package job

import (
	"GitlabCeForcedApprovals/enum"
	"GitlabCeForcedApprovals/helper"
	gitlab "gitlab.com/gitlab-org/api/client-go"
	"log"
)

const BotCommentBody = "This comment was created automatically to prevent merging before the MR is approved. Don't resolve this comment, it will be resolved automatically once approved."

type MergeEventHandler struct {
	Event  *gitlab.MergeEvent
	Gitlab *gitlab.Client
}

var botId int

func (receiver *MergeEventHandler) Handle() chan bool {
	result := make(chan bool)

	go func() {
		defer close(result)

		switch receiver.Event.ObjectAttributes.Action {
		case enum.MergeActionOpen,
			enum.MergeActionReopen,
			enum.MergeActionApproved,
			enum.MergeActionUnapproved,
			enum.MergeActionApproval,
			enum.MergeActionUnapproval,
			enum.MergeActionUpdate:
			receiver.HandleEvent(result)
			break
		default:
			result <- true
			break
		}
	}()

	return result
}

func (receiver *MergeEventHandler) HandleEvent(result chan bool) {
	projectId := receiver.Event.ObjectAttributes.SourceProjectID
	mergeId := receiver.Event.ObjectAttributes.IID
	reviewers := receiver.Event.Reviewers

	approvals, _, err := receiver.Gitlab.MergeRequests.GetMergeRequestApprovals(
		projectId,
		mergeId,
	)

	if err != nil {
		log.Printf("Failed getting approvals for MR %d in project %d: %s\n", mergeId, projectId, err)
		result <- false
		return
	}

	myComment, err := receiver.FindComment(projectId, mergeId)
	if err != nil {
		log.Printf("Failed getting comments for MR %d in project %d: %s\n", mergeId, projectId, err)
		result <- false
		return
	}

	if receiver.IsApproved(reviewers, approvals) {
		if myComment == nil || myComment.Notes[0].Resolved {
			result <- true
			return
		}

		_, _, err = receiver.Gitlab.Discussions.ResolveMergeRequestDiscussion(projectId, mergeId, myComment.ID, &gitlab.ResolveMergeRequestDiscussionOptions{
			Resolved: gitlab.Ptr(true),
		})
		if err != nil {
			log.Println("Failed resolving discussion:", err)
			result <- false
			return
		}

		result <- true
		return
	} else {
		if myComment != nil && !myComment.Notes[0].Resolved {
			result <- true
			return
		}

		if myComment == nil {
			myComment, _, err = receiver.Gitlab.Discussions.CreateMergeRequestDiscussion(projectId, mergeId, &gitlab.CreateMergeRequestDiscussionOptions{
				Body: gitlab.Ptr(BotCommentBody),
			})
			if err != nil {
				log.Println("Failed creating discussion:", err)
				result <- false
				return
			}
		}

		if myComment.Notes[0].Resolved {
			_, _, err = receiver.Gitlab.Discussions.ResolveMergeRequestDiscussion(projectId, mergeId, myComment.ID, &gitlab.ResolveMergeRequestDiscussionOptions{
				Resolved: gitlab.Ptr(false),
			})
			if err != nil {
				log.Println("Failed unresolving discussion:", err)
				result <- false
				return
			}
		}

		result <- true
		return
	}
}

func (receiver *MergeEventHandler) FindComment(projectId int, mergeId int) (*gitlab.Discussion, error) {
	comments, _, err := receiver.Gitlab.Discussions.ListMergeRequestDiscussions(projectId, mergeId, &gitlab.ListMergeRequestDiscussionsOptions{
		PerPage: 1000,
	})

	if err != nil {
		return nil, err
	}

	if botId == 0 {
		me, _, err := receiver.Gitlab.Users.CurrentUser()
		if err != nil {
			return nil, err
		}
		botId = me.ID
	}

	for _, comment := range comments {
		for _, note := range comment.Notes {
			if note.Author.ID != botId {
				continue
			}

			return comment, nil
		}
	}

	return nil, nil
}

func (receiver *MergeEventHandler) IsApproved(reviewers []*gitlab.EventUser, approvals *gitlab.MergeRequestApprovals) bool {
	if reviewers == nil || len(reviewers) == 0 {
		return approvals.Approved
	}

	approvedByIds := helper.SliceMap(approvals.ApprovedBy, func(item *gitlab.MergeRequestApproverUser) int {
		return item.User.ID
	})
	reviewerIds := helper.SliceMap(reviewers, func(item *gitlab.EventUser) int {
		return item.ID
	})

	if len(helper.SliceIntersect(reviewerIds, approvedByIds)) != len(reviewerIds) {
		return false
	}

	return true
}
