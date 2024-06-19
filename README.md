# Liscremon
Liscremon, Linux Screen Monitor, is not the proposed name of this project yet, all intention is to make this available on [Windows](https://en.wikipedia.org/wiki/Microsoft_Windows) and [MacOS](https://en.wikipedia.org/wiki/MacOS), which are very doable.

It's Linux for now because I'm a Linux person, a developer, a Go developer 😎 and development is sweetiest on Linux (this claim is very debatable, MacOS is popular with developers buh I do not own a Mac).

This README.md is meant to invite contributions to the project, cos am working too slow on it or mostly probable because I just want more opinions on the features that would be useful (building in secret is not so much fun I think). So if you want to touch grass with this codebase, THANK YOU 🙌🏿🙌🏿🙌🏿 and here are the things you should know.

## Project Scope:
The suitable name for this project is something like a 'screen monitor' without any operating system name attached. The central feature is to monitor **desktop screentime and application usage** something like the iOS screentime in the iPhone Settings or projects like [ActivityWatch](https://activitywatch.net/), [RescueTime](https://www.rescuetime.com/), [ManicTime](https://www.manictime.com/) or [WakaTime](https://wakatime.com/).

The plan is to follow more on the path with ActivityWatch, but in [Go](https://go.dev/), ActivityWatch is built with [Python](https://www.python.org/), with plans to have some part of it now in [Rust](https://www.rust-lang.org/). 

So far on this project, the central feature is available;

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


As for me, this page, the 'Sceentime' from the left menu is done. I said as for me, meaning any new idea is welcome to add or teardown. The other features include;

> **Tasks**

The third option on the left menu, is meant to consist of a 'Reminder' and 'App Limit', Both of this feauture I plan to make availabe for the first developer release (hopefully).

The **Reminder system** is just something that sends a desktop notification, with or without sound, as created for a task by the user, with a two pre-notification, also with an option to either launch an app on *StartTime*. Everything here is up for debate but I believe these are sane defaults. 

[implementation details is discussed below]().


The **App Limit** is also a notification that just tells you when your usage for a particular application is reached, also with the option to exist the app on limit reached. Now, this one has so few options in mind where we can debate about;

- a limit like most stuff has it, it tells you, one time, that the limit has been reached for that day, like the iOS, you can maybe extend for another 5 or 10 minute..ish.

- the first one seems like a mobile phone stuff, this option is something that tells you repeatedly you reached that limit, say you have a 2hrs App limit on VLC, you'll get a notification 3 times that day if you used VLC for 7hrs.

The one I'll be calling the sane default and would implement, is the first one, with no such thing as extension like iOS, this one could be a one time thing or a everyday thing (recurring daily). For the recurring daily, the limit only tick off the first time i.e once in that day.
