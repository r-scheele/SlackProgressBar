# SlackProgressBar

SlackProgressBar is a Go library for creating a real-time progress bar in a Slack channel. It provides a way to visualize the progress of long-running tasks directly in Slack, so your team can stay updated on the task's progress.

## Features

- Real-time progress updates in Slack
- Customizable progress bar with ETA
- Pause, resume, and log additional information while processing
- Memory usage display

## Installation

```bash
go get github.com/r-scheele/SlackProgressBar@v0.1.0
```

## Usage

1. First, import the library in your Go file:

```go
import (
    "github.com/r-scheele/SlackProgressBar/bar"
)
```

2. Create a new `SlackProgress` instance with your Slack token and channel:

```go
sp := &bar.SlackProgress{
    Token:   "your-slack-token",
    Channel: "#your-channel",
}
```

3. Define your list of items and the function to process each item:

```go
items := make([]string, 100)
for i := 1; i <= 100; i++ {
    items[i-1] = fmt.Sprintf("item%d", i)
}

func processItem(item string) {
    // Your item processing code here
}
```

4. Call `Iter` on your `SlackProgress` instance to process the items and update the progress bar in Slack:

```go
sp.Iter(items, processItem)
```

## Configuration

You can customize the appearance and behavior of the progress bar by setting additional fields on the `SlackProgress` struct:

```go
sp := &SlackProgress{
    Token:           "your-slack-token",
    Channel:         "#your-channel",
    Suffix:          "%",
    ProgressBarChar: "▓",
    Precision:       1,
    RateLimit:       time.Millisecond * 500,
    CompletionCallback: func() {
        fmt.Println("Task completed!")
    },
}
```



### Creating a Slack App

1. **Navigate to Slack API Page**:
   - Go to the [Slack API page](https://api.slack.com/apps) and click on `Create New App`.

2. **Name Your App**:
   - Enter a name for your app (e.g., `SlackProgressBar`) and select the Slack workspace where you want the app to reside, then click `Create App`.

3. **Configure Basic Information**:
   - Under the `Basic Information` section, you can set up various information about your app such as its name, description, and icons.

### Setting Up OAuth and Permissions

1. **Navigate to OAuth & Permissions**:
   - On the left sidebar, click on `OAuth & Permissions`.

2. **Add OAuth Scopes**:
   - Under the `Bot Token Scopes` section, click `Add an OAuth Scope` and add the following scopes:
     - `chat:write` (to send messages)
     - `chat:write.public` (if you want to post to channels without joining)

3. **Install App to Workspace**:
   - Click on `Install to Workspace` button at the top of the page.
   - You'll be redirected to a permission request screen in your Slack workspace. Click `Allow` to install the app and grant it the requested permissions.

4. **Copy the Bot User OAuth Token**:
   - Once the app is installed, you'll be redirected back to the `OAuth & Permissions` page.
   - Under the `OAuth Tokens for Your Workspace` section, you'll find a `Bot User OAuth Token`. This token will start with `xoxb-`. Copy this token as you'll need it for your Go code.

### Using the Token in Your Go Code

Now that you have the OAuth token, you can use it in your Go code as follows:

```go
sp := &SlackProgress{
    Token:           "xoxb-your-token-here",
    Channel:         "#your-channel",
    // ... other configuration ...
}
```

Replace `"xoxb-your-token-here"` with the actual Bot User OAuth Token you copied earlier.

### Running Your Go Code

Now you can run your Go code, and it should be able to post and update messages in the specified Slack channel using the credentials of the Slack app you created.

This setup allows your Go code to interact with the Slack API and post progress updates to your Slack workspace.

### Note:

- Make sure that the channel you specify in your Go code exists in your Slack workspace.
- If you need to post messages to multiple channels or private channels, you may need to invite the Slack app bot user to those channels within Slack.

## Contributing

We welcome contributions! Please feel free to submit a pull request with any improvements.

## License

This project is licensed under the MIT License.
```

You can copy and paste this template into a `README.md` file in the root of your GitHub repository, and then customize it to match the exact features and usage of your library.