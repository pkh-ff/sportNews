# SportNews

***api***

API服務

***crawler***

爬蟲服務，透過常駐程序的方式，定期去爬取資料，預設爬取時間區間如下

新聞爬蟲：爬取[https://sports.ndtv.com/cricket/news](https://sports.ndtv.com/cricket/news)資料，預設每2小時執行一次。

爬行榜爬蟲：爬取[https://www.bcci.tv/international/men/rankings](https://www.bcci.tv/international/men/rankings)資料，預設每12小時執行一次。

新聞圖片同步：將存在DB中的新聞封面進行轉存動作，預設每30分鐘執行一次。

> 1. 爬取時間區間可以透過[服務參數](#服務參數)中的`PROCESS_NEWS`與`PROCESS_RANKING`分別進行設定
> 2. 下次執行時間是根據依上次執行時間來決定

## Development Guide

### 開發環境需求

- Golang 1.24
- MySQL 8.0
- Docker（可選）
- Make（可選）

### 前置作業

***API與爬蟲服務運行參數***

將`conf`目錄下的yaml設定複製到專案根目錄並重新命名

```shell
# API Server設定檔
cp ./conf/app.conf.example.yaml ./app.conf.yaml

# 爬蟲服務設定檔
cp ./conf/process.conf.example.yaml ./process.conf.yaml 
```

依據[服務參數](#服務參數)中的說明進行相關參數設定

***makefile用環境變數***

複製專案根目錄下的`.env.example`，並改名為`.env`

```shell
cp ./.env.example .env
```

依據[Makefile 環境變數](#Makefile-環境變數)中的說明，進行相關環境變數設定

> 可以根據有沒有要使用make來決定要不要進行該步驟設置

### db migration

step1. 請先建立`sport_news` DB

```sql
CREATE DATABASE IF NOT EXISTS `sport_news` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci
```

step2. 執行db migration

以下方法擇一即可

方法 1:
將下面目錄`migration\sql`目錄下的SQL，在剛剛建立的DB下執行

方法 2: 執行下面指令

```shell
 make migrations
```

### build & run

***api***

執行下面指令，來運行API Server

```shell
go run .\cmd\app\main.go
```

***crawler***

執行下面指令，來運行爬蟲服務

```shell
go run .\cmd\crawler\main.go
```

## 部署指南

step1. build

建構docker image

make build SERVICE=`<服務名稱>` VERSION=`<服務版號>`

```shell
# 建構 http-server 0.0.1版本的Docker Image
make build SERVICE=http-server VERSION=0.0.1
# 建構 crawler-server 0.0.1版本的Docker Image
make build SERVICE=crawler-server VERSION=0.0.1
```

推送docker image

make push SERVICE=`<服務名稱>` VERSION=`<服務版號>`

```shell
# 推送 Ver 0.0.1的api-server Docker Image
make push SERVICE=http-server VERSION=0.0.1
# 推送 Ver 0.0.1的crawler-server Docker Image
make push SERVICE=crawler-server VERSION=0.0.1
```

step2. deploy

make run SERVICE=`<服務名稱>` VERSION=`<服務版號>`

運行 API服務

```shell
make run SERVICE=http-server VERSION=0.0.1 
```
運行爬蟲服務

```shell
make run SERVICE=crawler-server VERSION=0.0.1 
```

> 如果不使用make去執行，可以從[`makefile`](./makefile)檔案對照要執行命令去查看相關指令

## API Doc

使用postman將下面三個檔案[dev_environment.json](./docs/api/sport-news-dev.postman_environment.json)、[local_environment.json](./docs/api/sport-news-local.postman_environment.json)、[prod_environment.json](./docs/api/sport-news-prod.postman_environment.json)、[api_collection.json](./docs/api/SportNews.postman_collection.json)匯入，然後再根據要使用需求去指定postman使用環境變數設定黨。

## Other

### 環境變數&服務參數

#### 服務參數

***API Server***

| 變數 / 參數            | 說明                            |
|--------------------|-------------------------------|
| APP_NAME           | 服務名稱                          |
| APP_ADDR           | 服務位址(port號)                   |
| APP_DEBUG          | 是否啟用debug模式                   |
| APP_READ_TIMEOUT   | Http Request Read timeout 時效  |
| APP_WRITE_TIMEOUT  | Http Request Write timeout 時效 |
| DB_MASTER_HOST     | 主資料庫 連線 Host address          |
| DB_MASTER_USERNAME | 主資料庫 帳號                       |
| DB_MASTER_PASSWORD | 主資料庫 密碼                       |
| DB_MASTER_PORT     | 主資料庫 連線Port號                  |
| DB_SLAVER_HOST     | 從資料庫 連線 Host address          |
| DB_SLAVE_USERNAME  | 從資料庫 帳號                       |
| DB_SLAVE_PASSWORD  | 從資料庫 連線Port號                  |
| DB_SLAVE_PORT      | 從資料庫 連線Port號                  |

***Crawler Server***

| 參數                  | 說明                   |
|---------------------|----------------------|
| APP_NAME            | 服務名稱                 |
| APP_DEBUG           | 是否啟用debug模式          |
| APP_PROCESS_NEWS    | 新聞爬蟲執行間格時間(秒)        |
| APP_PROCESS_RANKING | 爬行榜爬蟲執行間格時間(秒)       |
| APP_PROCESS_PICTURE | 同步新聞封面圖片間格時間(秒)      |
| DB_MASTER_HOST      | 主資料庫 連線 Host address |
| DB_MASTER_USERNAME  | 主資料庫 帳號              |
| DB_MASTER_PASSWORD  | 主資料庫 密碼              |
| DB_MASTER_PORT      | 主資料庫 連線Port號         |
| DB_SLAVER_HOST      | 從資料庫 連線 Host address |
| DB_SLAVE_USERNAME   | 從資料庫 帳號              |
| DB_SLAVE_PASSWORD   | 從資料庫 連線Port號         |
| DB_SLAVE_PORT       | 從資料庫 連線Port號         |

> `服務參數`可以使用yaml或環境變數來進行設定，並以環境變數為優先

#### Makefile 環境變數

| 變數              | 說明                                 |
|-----------------|------------------------------------|
| PROJECT_PATH    | migration用sql目錄                    |
| DOCKER_REGISTRY | docker image repository's host     |
| DOCKER_USER     | docker image repository's username |
| DOCKER_PASSWORD | docker image repository's password |
| DB_HOST         | 資料庫連線Host                          |
| DB_USERNAME     | 資料庫連線帳號                            |
| DB_PASSWORD     | 資料庫連線密碼                            |

> 只支援透過環境變數方式設定

### db Schema

***news***

| 欄位名稱         | 類型           | 註解           |
|--------------|--------------|--------------|
| id           | int          | 新聞 ID        |
| title        | varchar(255) | 標題           |
| description  | varchar(255) | 描述           |
| cover        | varchar(155) | 封面           |
| cover_source | varchar(100) | 原始封面連結       |
| link         | varchar(100) | 新聞連結         |
| content      | text         | 內文           |
| status       | tinyint      | 狀態，上架:1、下架:0 |
| source       | varchar(30)  | 來源           |
| pub_date     | timestamp    | 發佈時間         |
| create_at    | timestamp    | 建立時間         |
| update_at    | timestamp    | 更新時間         |

***video***

| 欄位名稱        | 類型           | 註解           |
|-------------|--------------|--------------|
| id          | int          | 視頻 ID        |
| title       | varchar(100) | 標題           |
| description | varchar(255) | 描述           |
| cover       | varchar(255) | 封面           |
| link        | varchar(255) | 連結           |
| status      | tinyint      | 狀態，上架:1、下架:0 |
| create_at   | timestamp    | 建立時間         |
| update_at   | timestamp    | 更新時間         |

***rank***

| 欄位名稱      | 類型          | 註解    |
|-----------|-------------|-------|
| id        | int         | 排名 ID |
| type      | varchar(10) | 排名類型  |
| data      | text        | 排名資料  |
| date      | date        | 資料日期  |
| create_at | timestamp   | 建立時間  |

### crawler status

***News***

- [ ] [Cricbuzz](https://www.espncricinfo.com/cricket-news)
- [ ] [ESPNcricinfo](https://www.espncricinfo.com/)
- [ ] [BCCI官方网站](https://www.bcci.tv/international/men/news)
- [x] [NDTV Sports - Cricket](https://sports.ndtv.com/cricket/news)
- [ ] [Sportskeeda](https://www.sportskeeda.com/cricket)

***Ranking***

- [ ] [Cricbuzz](https://www.cricbuzz.com/cricket-stats/icc-rankings/men/teams)
- [x] [BCCI官方网站](https://www.bcci.tv/international/men/rankings)
- [ ] [Sportskeeda](https://www.sportskeeda.com/cricket/icc-rankings)
