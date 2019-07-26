# RSS Search API

[![Build Status](https://travis-ci.org/shhj1998/rss-search-api.svg?branch=master)](https://travis-ci.org/shhj1998/rss-search-api)
[![Coverage Status](https://coveralls.io/repos/github/shhj1998/rss-search-api/badge.svg?branch=master)](https://coveralls.io/github/shhj1998/rss-search-api?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/shhj1998/rss-search-api)](https://goreportcard.com/report/github.com/shhj1998/rss-search-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![](https://godoc.org/github.com/shhj1998/rss-search-api/rsserver?status.svg)](https://godoc.org/github.com/shhj1998/rss-search-api/rsserver)

This repository is a server code for fetching and searching RSS data efficiently. There are also many useful Go packages to help you develop your own RSS project! Our version is automatically released to <https://rss-search-api.herokuapp.com>. We are running the server with Daum RSS. You can try GET apis for example.

## Table of Contents

1. [Packages](#packages)
2. [Getting Start](#getting-start)
3. [API](#api)
   - [Item API](#item-api)
   - [Channel API](#channel-api)
4. [Schema](#schema)
   - [ER Model](#er-model)
   - [Table](#table)
5. [Performance](#performance)
   - [Test DB Info](#test-db-info)
   - [Results](#results)
6. [Library Used](#library-used)
7. [See Also](#see-also)

## Packages

- [rsserver](https://github.com/shhj1998/rss-search-api/tree/master/rsserver) - Our main package that provides many features to act with our rss database. More information is available in our [documentation](https://godoc.org/github.com/shhj1998/rss-search-api/rsserver).
- [rssapi](https://github.com/shhj1998/rss-search-api/tree/master/rssapi) - A simple REST api server that uses rsserver package. Available apis are described below.

## Getting Start

Assuming you already have a recent version of Go installed, pull down the code with `go get`.
```bash
go get github.com/shhj1998/rss-search-api
```
Go into the source directory and pull down the project dependencies by using [`dep`](https://github.com/golang/dep):
```bash
cd $GOPATH/src/github.com/shhj1998/rss-search-api
dep ensure
```
Now, generate `.env` file in the same directory. It is used to configure database connection. Sample .env file is [here](https://github.com/shhj1998/rss-search-api/blob/master/.env.example).
```bash
vim .env
```

Now, run
```bash
go run main.go
```

You will see the following sentences.
```bash
INFO : 2019/07/26 23:08:25.163777 main.go:34: Total CPUs : 4
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/channel           --> github.com/shhj1998/rss-search-api/rssapi.(*Server).SelectChannels-fm (3 handlers)
[GIN-debug] GET    /api/v1/channel/items/    --> github.com/shhj1998/rss-search-api/rssapi.(*Server).SelectChannelsWithItems-fm (3 handlers)
[GIN-debug] GET    /api/v1/channel/items/count/:count --> github.com/shhj1998/rss-search-api/rssapi.(*Server).SelectChannelsWithItems-fm (3 handlers)
[GIN-debug] POST   /api/v1/channel           --> github.com/shhj1998/rss-search-api/rssapi.(*Server).CreateChannel-fm (3 handlers)
[GIN-debug] GET    /api/v1/item              --> github.com/shhj1998/rss-search-api/rssapi.(*Server).SelectItems-fm (3 handlers)
[GIN-debug] Listening and serving HTTP on :80
```

## API

### Item API

These apis help you to fetch and create rss item information.

#### GET /api/v1/item

#### Request

```bash
curl -i -H 'Accept: application/json' http://localhost/api/v1/item
```

#### Response

```json
[
  {
    "title": "...",
    "description": "...",
    "link": "...",
    "published": "...",
    "author": {
      "name": "..."
    },
    "guid": "...",
    "enclosures": [
      {
        "url": "...",
        "type": "..."
      },
      ...
    ]
  }, 
  ...
]
```

### Channel API

These apis help you to fetch and create rss channel information.

#### GET /api/v1/channel

#### Request

```bash
curl -i -H 'Accept: application/json' http://localhost/api/v1/channel
```

#### Response

```json
[
  {
    "id": 1,
    "rsslink": "https://media.daum.net/syndication/society.rss",
    "title": "다음뉴스 - 사회Top RSS",
    "link": "http://media.daum.net/society/",
    "items": null,
    "feedType": "",
    "feedVersion": ""
  },
  {
    "id": 9,
    "rsslink": "https://media.daum.net/syndication/entertain.rss",
    "items": null,
    "feedType": "",
    "feedVersion": ""
  }
]
```

#### POST /api/v1/channel

#### Request

```bash
curl -i -H 'Accept: application/json' -X POST \ 
	--data '{"rss": "https://media.daum.net/syndication/entertain.rss"}' \
	http://localhost/api/v1/channel
```

#### Response

```json
{
  "success": "successfully created a new Channel"
}
```

#### Error

- If you forgot adding request body with rss link, 

```json
{
  "error": "Request body doesn't match the api... Please read our api docs"
}
```

- If your rss link is invalid,

```json
{
  "error": "http error: 400 Bad Request"
}
```

- If your rss link is not a rss,

```json
{
  "error": "Failed to detect feed type"
}
```

- If you are creating an already existing rss channel,

```json
{
  "error": "Error 1062: Duplicate entry '...' for key 'rss_link'"
}
```

#### GET /api/v1/channel/items

#### Request - fetch all channels with items

```bash
curl -i -H 'Accept: application/json' http://localhost/api/v1/channel/items
```

#### Response

```json
[
  {
    "id": 1,
    "rsslink": "https://media.daum.net/syndication/society.rss",
    "title": "...",
    "link": "...",
    "items": [...],
    "feedType": "",
    "feedVersion": ""
  },
  {
    "id": 9,
    "rsslink": "https://media.daum.net/syndication/entertain.rss",
    "items": null,
    "feedType": "",
    "feedVersion": ""
  }
]
```

#### Request - fetch specific channels with items

You can filter the channel you only want by query strings. It is available to put multiple channel's id.

```bash
curl -i -H 'Accept: application/json' http://localhost/api/v1/channel/items?id=1&id=2000
```

#### Response

```json
[
  {
    "id": 1,
    "rsslink": "https://media.daum.net/syndication/society.rss",
    "title": "...",
    "link": "...",
    "items": [...],
    "feedType": "",
    "feedVersion": ""
  }
]
```

#### GET /api/v1/channel/items/count/{count}

You can also use the feature of channel filtering with query string described in **GET /api/v1/channel/items**.

### Request

```bash
curl -i -H 'Accept: application/json' http://localhost/api/v1/channel/items/count/1
```

#### Response

```json
[
  {
    "id": 1,
    "rsslink": "https://media.daum.net/syndication/society.rss",
    "title": "...",
    "link": "...",
    "items": [...],
    "feedType": "",
    "feedVersion": ""
  },
  {
    "id": 9,
    "rsslink": "https://media.daum.net/syndication/entertain.rss",
    "items": null,
    "feedType": "",
    "feedVersion": ""
  }
]
```

## Schema

### ER Model

- Entity
  - Channel : RSS Channels provide items(news). It doesn't have a specific primary key, so we set a new feature 'channel_id' as a primary key.
  - Item : Items are the news published by RSS channels. In RSS 2.0, their is a primary key called 'guid'.
  - Enclosure : Enclosures are media that is used in items like image, video. Usually, one url only have one media, so we set url as a primary key.
- Relationship
  - Publish : This relationship represents what items were published by a channel. It is many-to-many.
  - Enclosure : This relationship represents which enclosures are used in a item. It is one-to-many.

![Imgur](https://i.imgur.com/IOEGaAR.png)

### Table

You can see the tables schema [here](https://github.com/shhj1998/rss-search-api/blob/master/rsserver/model.go). We add some b-tree index to higher the performance. If you are interested what is b-tree index, you can look [here](https://www.quora.com/What-is-B-trees-index-in-SQL).

- Item
  - Because 'guid' is a long string, it will lower performance in search. So we made a new primary key which is integer named 'item_id'. Instead, we add a unique contraint to guid.
  - We add a b-tree index on item_id because it will be frequently used in joining with Publish table.
- Channel
  - Nothing special, just add a auto_increment integer primary key 'channel_id' with a b-tree index.
- Publish
  - To satisfy [BCNF](https://en.wikipedia.org/wiki/Boyce–Codd_normal_form) constraints - lower redundancy, we change Publish relationship to a table instead adding a feature in Item. 
  - Because it is many-to-many, (item, channel) is a primary key. Each of them is a foreign key.
  - Add a b-tree index to 'channel'.
- Enclosure
  - Because Media is a one-to-many relationship, it don't need a additional table to lower redundancy.
  - Add a b-tree index and foreign key constraint to 'item'.

## Performance
### Test DB Info

We used 8 rss services and 3594 rss items(news) to test our api performance. The details of the rss services are described below. The data we used to test are available [here](https://drive.google.com/file/d/1oYRqRyp56PQuHcnNJW9Non5rwjnasSLL/view?usp=sharing).

| Name                           | Link                                                | Rows |
| ------------------------------ | --------------------------------------------------- | ---- |
| 다음뉴스 - 사회Top RSS         | https://media.daum.net/syndication/society.rss      | 387  |
| 다음뉴스 - 시사 주요뉴스 RSS   | https://media.daum.net/syndication/entertain.rss    | 975  |
| 다음뉴스 - 스포츠 주요뉴스 RSS | https://media.daum.net/syndication/today_sports.rss | 810  |
| 다음뉴스 - 정치Top RSS         | https://media.daum.net/syndication/politics.rss     | 312  |
| 다음뉴스 - 경제Top RSS         | https://media.daum.net/syndication/economic.rss     | 431  |
| 다음뉴스 - 국제Top RSS         | https://media.daum.net/syndication/foreign.rss      | 302  |
| 다음뉴스 - 문화/생활Top RSS    | https://media.daum.net/syndication/culture.rss      | 202  |
| 다음뉴스 - Tech Top RSS        | https://media.daum.net/syndication/digital.rss      | 163  |

### Results

Details of the result are described below. The code used to test performance is [here](https://github.com/shhj1998/rss-search-api/blob/master/rssapi/channel_test.go). It will not work in your local repository because it must be with a .env file that contains the test db information. We tested three api that were most likely to be used. **GET /api/v1/channel/items/count/{count}** apis show almost same performance although the count value changes.

| API                                | Time(ms)   |
| ---------------------------------- | ---------- |
| GET /api/v1/channel                | 39.837295  |
| GET /api/v1/channel/items          | 298.150581 |
| GET /api/v1/channel/items?id=3     | 61.691366  |
| GET /api/v1/channel/items/count/1  | 277.327405 |
| GET /api/v1/channel/items/count/2  | 273.187845 |
| GET /api/v1/channel/items/count/3  | 281.879731 |
| GET /api/v1/channel/items/count/4  | 347.423070 |
| GET /api/v1/channel/items/count/5  | 283.924582 |
| GET /api/v1/channel/items/count/6  | 274.330836 |
| GET /api/v1/channel/items/count/7  | 279.774116 |
| GET /api/v1/channel/items/count/8  | 270.644527 |
| GET /api/v1/channel/items/count/9  | 339.821699 |
| GET /api/v1/channel/items/count/10 | 345.249172 |

![Imgur](https://i.imgur.com/ztFsvEJ.png)

## Library Used

- [gofeed](github.com/mmcdole/gofeed) - Provides rss parser and related types.
- [sql-mock](github.com/DATA-DOG/go-sqlmock) - Used to mock sql database and connection for testing.
- [testify](github.com/stretchr/testify) - Provides useful test functions like assertions.
- [logger](github.com/google/logger) - Custom logger library.

## See Also

- [RSS 2.0](https://cyber.harvard.edu/rss/rss.html) - More information about RSS specification.
- [android-rss-viewer](https://github.com/Park-Wonbin/android-rss-viewer) - Android application that uses rss-search-api.
