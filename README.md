# Bronco Build Tracker
When Ford is building a vehicle that is a special order
there is a tracking website to show the progress and provide
access to the window sticker when it is ready. Instead
of refreshing the page on every opportunity, I wrote some
Bash scripts to notify me of changes. To share this capability and
take an opportunity to play with Go, I created this. Pleas forgive
the overall quality as it is more of a toy than anything. It might be
useful, but it comes with no guarantees.

The app will:
- Check the status site and track progress (pizza tracker)
- Check the window sticker link and track changes
- Check the authenticated status site and track progress (backdoor tracker)
- Notify you of changes via Textbelt or Pushover
- Run the checks on a regular interval

There will be two files created in the same directory as the application:
- settings.json - This will contain all the application settings. It will contain personal info
so do not share it with anybody.
- state.json -  This will hold information about the state of your build and is safe to share.

The program will prompt for settings or you can provide them
with CLI options.
```
Application Options:
  -o, --order-number=   Order number (required for pizza tracker)
  -v, --vin=            Vehicle VIN (required for pizza tracker)
  -r, --reservation=    Reservation ID (required for backdoor tracker)
      --refreshToken=   Refresh token (required for backdoor tracker)
  -i, --interval=       Interval between checks
  -p, --phone=          Phone number to text (required for textbelt
                        notification)
      --textbelt-key=   Textbelt API key (required for textbelt notification)
      --pushover-token= Pushover token (required for pushover notification)
      --pushover-user=  Pushover user (required for pushover notification)
      --once            Flag to run check once and stop.
  -q, --quiet           Quiet run - don't prompt for any settings

Help Options:
  -h, --help            Show this help message
```
### Textbelt
Textbelt provides a way to send SMS messages from apps like this.
If you would like to receive SMS notifications you need to setup textbelt.
You do have to pay for the notifications, but it is cheap. 
For Textbelt notifications, setup a key and pay for text messages
at http://textbelt.com. You can get 50 texts for $3

### Pushover
For Pushover notifications, download the app and go to http://pushover.net
and setup a token and user. The app is $5 one time fee (free for 1 month)
and 10,000 notification per month for free.

### Getting a refresh token
A refresh token is part of the OAuth authentication method used to login to 
the Ford website. This token has a long lifespan and will allow us to grab 
the info we need. Please note, never share your token with anybody else as this
could compromise your account! Getting this token might prove difficult if you are not familiar with
this kind of thing, but it is interesting because it is very reliable and exposes the vehicleStatusCode
that seems to change more often than the visual components of the two tracker...although I don't know 
what the values mean yet.
To get the token:
- Go to https://www.ford.com/buy/manage.html?reservationId=xxxxxx (using your reservation ID) in Chrome
- Right-click on the page and select "Inspect"
- Login
- Click on the Network tab in the pane that opened up when you clicked "Inspect"
- In the filter in the top left section of the pane type "api.mps.ford.com"
- In the list of requests, there will be two with the name "token". Click on the second one and look at 
the Preview sub-pane
- Find the value of the refresh_token field and you will have the value