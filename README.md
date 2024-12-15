Simulate approval requirements for GitLab Community Edition. A bot user creates an unresolved thread for every merge 
request and resolves it automatically when the MR is approved.

## Installing

You can deploy the app using docker, using serverless (CloudFormation) or just build the binary yourself and run it 
anywhere. In any case, two environment variables are required:

- `GITLAB_ACCESS_TOKEN` - the access token of the bot user who will post and resolve the blocking thread 
  (the bot needs at least Developer access to be able to resolve conversations)
- `GITLAB_BASE_URL` - the url (including protocol) to your GitLab instance

Afterwards, you need to go to your instance and set the webhook with the deployed url. You can either configure it for 
the whole instance in the admin, or per each project.

### Serverless

- `export GITLAB_ACCESS_TOKEN=your-bot-access-token`
- `export GITLAB_BASE_URL=https://gitlab.example.com`
- `yarn install` (or `npm install`)
- `serverless deploy --stage prod --verbose`
- The output of the command should list `HttpApiUrl` that looks like `https://xxxxxxxxxx.execute-api.eu-central-1.amazonaws.com`
- Put the above url with `/webhooks` added as a path as your webhook URL (for example `https://xxxxxxxxxx.execute-api.eu-central-1.amazonaws.com/webhooks`)
