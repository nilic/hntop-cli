# hntop-cli 
CLI utility for displaying top [Hacker News](https://news.ycombinator.com/) articles for a given time period. Posts are sorted based on number of points, then number of comments.

Relies on [HN Search API](https://hn.algolia.com/api). Please bear in mind that Algolia limits the number of API requests from a single IP to 10,000 per hour.

## Installation

Just grab the archive for your OS and platform from the [Releases](https://github.com/nilic/hntop-cli/releases) page and extract it somewhere. Optionally, you can add `hntop` to your [$PATH](https://gist.github.com/nex3/c395b2f8fd4b02068be37c961301caa7).

`hntop` is also available as a [container image](https://github.com/nilic/hntop-cli/pkgs/container/hntop-cli).

## Usage

```
$ hntop -h
NAME:
   hntop - display top Hacker News posts

USAGE:
   hntop [global options] command [command options] [arguments...]

VERSION:
   v9.9.9

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

   Mail options:

   --mail-auth value      Mail server authentication mechanism, one of: [login plain crammd5 xoauth2 none]. (default: login) [$HNTOP_MAIL_AUTH]
   --mail-from value      Mail From address. [$HNTOP_MAIL_FROM]
   --mail-password value  Mail server password. [$HNTOP_MAIL_PASSWORD]
   --mail-port value      Mail server port. (default: 587) [$HNTOP_MAIL_PORT]
   --mail-server value    Mail server. [$HNTOP_MAIL_SERVER]
   --mail-tls value       Mail server TLS policy, one of: [mandatory opportunistic notls]. (default: mandatory) [$HNTOP_MAIL_TLS]
   --mail-to value        Mail To address. [$HNTOP_MAIL_TO]
   --mail-username value  Mail server username. [$HNTOP_MAIL_USERNAME]

   Miscellaneous:

   --output value, -o value  Output format, one of: [list mail]. (default: "list") [$HNTOP_OUTPUT]

   Search options:

   --count value, -c value  Number of results to retrieve, must be between 1 and 1000. (default: 20) [$HNTOP_COUNT]
   --front-page, -f         Display current front page posts. If selected, all other flags are ignored. (default: false) [$HNTOP_FRONT_PAGE]
   --tags value, -t value   Filter results by post tag. Available tags: [story poll show_hn ask_hn]. (default: "story,poll,show_hn,ask_hn") [$HNTOP_TAGS]

   Time interval:

   --from value            Start of the time range to show top HN posts from in RFC3339 format. Used in conjuction with --to. [$HNTOP_FROM]
   --last value, -l value  Interval since current time to show top HN posts from, eg. "12h" (last 12 hours), "6m" (last 6 months). [$HNTOP_LAST]
   --to value              End of the time range to show top HN posts from in RFC3339 format. Used in conjuction with --from. If omitted, current time will be used. [$HNTOP_TO]
```

## Examples

### Get top HN posts from last X hours, days, weeks, months or years

Interval to show posts from is defined as `<length><unit>` since current time, eg. `12h` for posts from last 12 hours or `6m` for posts from last 6 months.

Available units: `h` - hour, `d` - day, `w` - week, `m` - month, `y` - year.

```
# get top HN posts from last week
hntop
hntop -l 1w

# get top HN posts from last 3 days
hntop -l 3d

# get top HN posts from last 9 months
hntop -l 9m

# get top HN posts from last 50 years
hntop -l 50y
```

### Get top HN posts in a custom timerange

Custom timerange can be defined using the RFC3339 format, ie. `yyyy-MM-dd'T'HH:mm:ss'Z'` for UTC or `yyyy-MM-dd'T'HH:mm:ss±hh:mm` for a specific timezone, where `±hh:mm` is the offset to UTC.

Examples: `2006-01-02T15:04:05Z` is 2 Jan 2006 15:40:05 UTC, while `2017-10-12T20:05:09+01:00` is 12 Oct 2017 20:05:09 CET.

```
# get top HN posts from 1 Jan 2018 to 1 July 2018
hntop --from 2018-01-01T00:00:00Z --to 2018-07-01T00:00:00Z

# get top HN posts from 24 Feb 2023 to now
hntop --from 2023-02-24T00:00:00Z

# get top HN posts from 24 Sep 2016 at 10 AM CET to 24 Sep 2016 at 12 AM CET
hntop --from 2016-09-24T10:00:00+01:00 --to 2016-09-24T12:00:00+01:00
```

### Misc

```
# get top posts currently on the HN front page
# posts will appear in different order than on front page because of the sorting algorithm
# all other flags are ignored in this case
hntop -f

# get 100 top HN posts from last week instead of default 20
hntop -c 100

# get top "Show HN" posts from last year
hntop -l 1y -t show_hn
```

### Output

Results can be printed out as a list (default) or sent by e-mail.

```
# output to console as a list (default)
hntop -o list

# send results by e-mail
hntop -o mail --mail-from bob@bob.com --mail-to alice@alice.com --mail-server smtp.bob.com --mail-port 587 --mail-tls mandatory --mail-auth login --mail-username bob --mail-password B0bsP@ss 
```
