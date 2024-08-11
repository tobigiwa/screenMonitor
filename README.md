# SCREEN MONITOR

## Liscremon

Liscremon, Linux Screen Monitor, is not the proposed name of this project yet, all intention is to make this available on [Windows](https://en.wikipedia.org/wiki/Microsoft_Windows) and [MacOS](https://en.wikipedia.org/wiki/MacOS), which are very doable.

It's Linux for now, because it is the development environment I use. I own a banged-up HP Pavilion, so Windows support is definately coming(WinScreMon). I can't say much of MacOS(DaScreMon), shit is too expensive.

This README.md is meant to invite contributions to the project, I believe am working too slow on it, mostly probable because I just want more opinions on the features that would be useful (building in secret is not so much fun I think). So if you want to touch grass with this codebase, THANK YOU 🙌🏿🙌🏿🙌🏿 and here are the things you should know.

## Project Scope:

The suitable name for this project is something like a 'screen monitor' without any operating system name attached. The central feature is to monitor **desktop screentime and application usage** something like the iOS screentime in the iPhone Settings or projects like [ActivityWatch](https://activitywatch.net/), [RescueTime](https://www.rescuetime.com/), [ManicTime](https://www.manictime.com/) or [WakaTime](https://wakatime.com/).

The plan is to follow more on the path with ActivityWatch, but in [Go](https://go.dev/), ActivityWatch is built with [Python](https://www.python.org/), with plans to have some part of it now in [Rust](https://www.rust-lang.org/).

So far on this project, the central feature is available;

> **Screentime**

<details>
    <summary>
        Weekly Screentime
    </summary>
    <img src="./images/weekly.png" alt="Weekly Report">
</details>

<details>
    <summary>
        Daily Screentime, on clicking on any of the weekly bar
    </summary>
    <img src="./images/daily.png" alt="Weekly Report">
</details>

<details>
    <summary>
        Application Screentime for the week
    </summary>
    <img src="./images/weekly-app.png" alt="Weekly Report">
</details>

As for me, this page, the 'Sceentime', is completed. I said as for me, meaning any new idea is welcome to add or teardown. The other features as noticed from the menu bar includes;

[implementation details is discussed below]().

> **Tasks**

The third option on the left menu, is meant to consist of a 'Reminder' and 'Daily App Limit', Both of this feauture are availabe for the first developer release (functional but the UI is very shitty). Sounds cool but not central to the scope of the project.

The **Reminder system** is just something that sends a desktop notification, with or without sound, as created for a task by the user, with a two pre-notification, also with an option to either launch an app on _StartTime_. Everything here is up for debate but I believe these are sane defaults.

[implementation details is discussed below]()

The **Daily App Limit** is also a notification that just tells you when your usage for a particular application is reached for a day, also with the option to exist the app on limit reached.

[implementation details is discussed below]().

> **Analytics**

This feature, unlike _Tasks_, is central to the application, it like a secondary feature to _Screentime_. This is where you can do advance analytics on your screentime data aside "weekly". A month screentime, a 3 week application screentime, in Line plot or whatever is fine presentation for such data. If you take a closer look at the "Weekly Screentime" image shared above, you would notice under the name of some application, a categorization exist, we can have analytics that make use of this categorization. Something like, a two week usage analytics with application categorized under "Entertainment and Gaming".

This part should be fun, ActivityWatch has something like that, it is called "Query" or something. On the discussion regarding Application Category, [see here](). This part, in my intention might not be available on the first developer release. Opinions on how this should be implemented and presented is what am looking forward too.

> **ToDo**

This one, not central at all, is the last thing we should handle if we agree it is "Okay" to have. The idea is to have a UI like the mobile app version of [Trello](https://trello.com/), where we have three columns; "ToDo", "Doing" and "Done"; with drag and drop. Ideally Todo items should be limited to two weeks. Seeing we have a Reminder system already, it would only be frontend heavy, so small work on the Go side of things.


## Project Architecture
By the folder structure of the codebase;
- cli: this is the main entrypoint to the **daemon service**. It has few commands to launch and stop the  it does the screen monitoring, task scheduling and database management. It is always running, so it has an autostart script. It is the **smDaemon binary**.
Both the cli/ and daemon/ makes up this

- 

-agent: This is the backend that talks to the daemon service, and also contains the frontend that displays to the user.
- TrayIcon: This is meant to be how the us
