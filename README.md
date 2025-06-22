# GeekNews Daily Bot

GeekNews RSS 피드를 크롤링하여 PostgreSQL 데이터베이스에 저장하고, Discord 웹훅을 통해 일일 뉴스 요약을 전송하는 봇입니다.

## 기능

- **자동 RSS 크롤링**: 매시간 59분에 GeekNews RSS 피드 크롤링
- **일일 Discord 알림**: 매일 00:00 UTC (09:00 KST)에 새로운 뉴스 Discord 전송
- **중복 처리**: URL 기반으로 중복 뉴스 방지 (UPSERT)
- **PostgreSQL 저장**: 모든 RSS 데이터를 구조화된 형태로 저장
- **전송 상태 추적**: 뉴스별 Discord 전송 상태 관리
- **환경변수 설정**: 데이터베이스 및 Discord 연결 정보 외부 설정
- **에러 로깅**: 크롤링 및 전송 실패 시 상세 로그 출력

## 설치 및 실행

### 1. 의존성 설치

```bash
go mod tidy
```

### 2. 데이터베이스 설정

PostgreSQL 데이터베이스에 스키마를 생성합니다:

```bash
psql -d your_database -f schema.sql
```

### 3. 환경변수 설정

다음 환경변수를 설정합니다:

| 변수명                   | 필수 | 기본값                                          | 설명                 |
|-----------------------|----|----------------------------------------------|--------------------|
| `DB_HOST`             | ✅  | -                                            | PostgreSQL 호스트     |
| `DB_PORT`             | ✅  | -                                            | PostgreSQL 포트      |
| `DB_USER`             | ✅  | -                                            | PostgreSQL 사용자명    |
| `DB_PASSWORD`         | ✅  | -                                            | PostgreSQL 비밀번호    |
| `DB_NAME`             | ✅  | -                                            | PostgreSQL 데이터베이스명 |
| `DB_SSLMODE`          | ❌  | `disable`                                    | PostgreSQL SSL 모드  |
| `RSS_FEED_URL`        | ❌  | `https://feeds.feedburner.com/geeknews-feed` | RSS 피드 URL         |
| `DISCORD_WEBHOOK_URL` | ✅  | -                                            | Discord 웹훅 URL     |

### 4. 빌드 및 실행

```bash
go build -o geek-news-bot main.go
./geek-news-bot
```

## 환경변수

## 데이터베이스 스키마

`news` 테이블에 다음 필드들이 저장됩니다:

- `id`: 뉴스 고유 ID (Primary Key, SERIAL)
- `url`: 뉴스 링크 (UNIQUE)
- `title`: 뉴스 제목
- `author`: 작성자 이름
- `content`: 뉴스 내용
- `published_at`: 뉴스 발행 시간
- `created_at`: 데이터베이스 저장 시간 (기본값: NOW())
- `sent`: Discord 전송 상태 (기본값: FALSE)

## 스케줄링

### RSS 크롤링

- **실행 주기**: 매시간 59분에 실행
- **Cron 표현식**: `59 * * * *`
- **목적**: GeekNews RSS 피드에서 새로운 뉴스 수집 및 데이터베이스 저장

### Discord 알림

- **실행 주기**: 매일 00:00 UTC (09:00 KST)에 실행
- **Cron 표현식**: `0 0 * * *`
- **목적**: 미전송 뉴스를 Discord로 전송 후 전송 상태 업데이트

## 로그

프로그램은 다음과 같은 로그를 출력합니다:

- 데이터베이스 연결 상태
- RSS 크롤링 결과 (성공/실패)
- 저장된 뉴스 항목 수 (신규/중복/에러)
- Discord 전송 결과 및 전송된 항목 수
- 스케줄링 상태 및 작업 ID

## 파일 구조

```
.
├── main.go                      # 메인 애플리케이션 엔트리포인트
├── schema.sql                   # PostgreSQL 스키마 정의
├── go.mod                      # Go 모듈 정의
├── go.sum                      # 의존성 체크섬
├── README.md                   # 프로젝트 문서
├── config/                     # 설정 관리
│   └── config.go
├── crawler/                    # RSS 크롤링
│   └── rss.go
├── database/                   # 데이터베이스 연결 및 저장
│   ├── connection.go
│   └── repository.go
├── discord/                    # Discord 웹훅 전송
│   └── webhook.go
├── models/                     # 데이터 모델 정의
│   ├── rss.go
│   └── database.go
└── scheduler/                  # 스케줄링 및 작업 조율
    └── scheduler.go
```

## 패키지 구조

### config

- 환경변수 로딩 및 설정 관리

### models

- RSS 피드 구조체 (`Feed`, `Entry`, `Author`, `Link` 등)
- 데이터베이스 구조체 (`News`)
- 데이터 변환 함수

### crawler

- RSS 피드 크롤링 기능
- HTTP 요청 및 XML 파싱

### database

- PostgreSQL 연결 관리
- 뉴스 엔트리 저장 로직 (중복 방지)
- 미전송 뉴스 조회 및 전송 상태 업데이트

### discord

- Discord 웹훅을 통한 뉴스 전송 기능
- 뉴스 목록을 마크다운 링크 형태로 변환

### scheduler

- gocron을 사용한 이중 스케줄링 (RSS 크롤링 + Discord 알림)
- 크롤링-저장-전송 파이프라인 조합

## Docker 실행

Docker를 사용하여 컨테이너로 실행할 수 있습니다:

```bash
# 이미지 빌드
docker build -t geek-news-bot .

# 환경변수와 함께 컨테이너 실행
docker run -d \
  -e DB_HOST=your_db_host \
  -e DB_PORT=5432 \
  -e DB_USER=your_username \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=geeknews_db \
  -e DISCORD_WEBHOOK_URL=your_webhook_url \
  --name geek-news-bot \
  geek-news-bot
```

## 주요 의존성

- **gocron/v2**: 스케줄링 라이브러리
- **pgx/v5**: PostgreSQL 드라이버
- **sqlx**: SQL 확장 라이브러리
- **godotenv**: 환경변수 로딩
