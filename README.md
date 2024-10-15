# Is there snow yet????

![](docs/nfty-notification.png)

## Usage

* CLI options:

```
Usage of itsy:
  -img-url string
        url of the image to download (mandatory)
  -name string
        name of the monitor (default "snow monitor")
  -nfty-callback-address string
        if set, you'll be redirected to this address when opening the notification
  -nfty-embed-image
        if set, it will embed the downloaded image to the notification (if size < 2Mo)
  -nfty-topic string
        nfty topic to send notifications when using nfty notifier
  -notifier string
        select notifier to use (term, ntfy) (default "term")
  -snow-only
        send notification only if snow has been detected
  -threshold float
        confidence threshold, in percent (100 = absolutely sure) (default 25)
```

* Usage example using a cronjob:

```
0 12 * * * /path/to/itsy -img-url "${IMG_URL?}" -name "${MONITOR_NAME?} (cron)" -notifier nfty -nfty-topic is-there-snow-yet  -threshold 25 -nfty-embed-image
```

## License

[MIT](https://choosealicense.com/licenses/mit/)