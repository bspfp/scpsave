# scpsave

유사 게임 클라우드 저장 with SCP

- [사용법](#사용법)
  - [clone](#clone)
  - [빌드](#빌드)
  - [설정 파일 만들기](#설정-파일-만들기)
  - [설정 파일 수정 후 이름 변경](#설정-파일-수정-후-이름-변경)
  - [실행](#실행)
- [설정 파일 내용](#설정-파일-내용)

## 사용법

### clone

```powershell
cd X:\Working\Folder
git clone https://github.com/bspfp/scpsave.git
```

### 빌드

in scpsave/cmd/scpsave

```powershell
go build -o scpsave.exe ./...
```

### 설정 파일 만들기

in scpsave/cmd/scpsave

```powershell
.\scpsave.exe -c
```

### 설정 파일 수정 후 이름 변경

in scpsave/cmd/scpsave

```powershell
Rename-Item -Path .\config.sample.yaml -NewName config.yaml
```

### 실행

```powershell
.\scpsave.exe
```

## 설정 파일 내용

```yaml
server_address: example.com:22
username: user
private_key_path: C:\Users\user\.ssh\id_rsa
remote_root: /remote/path
games:
  - name: Game1
    local_dir: C:\Users\user\Games\Game1
    file_patterns:
      - some.+\\.+\.save
      - .+\.dat
    program_name: game1.exe
  - name: Game2
    local_dir: C:\Users\user\Games\Game2
    file_patterns:
      - .+\.sav
      - .+\.dat
    program_name: C:\Game\Folder\game2.exe
```

| 항목                | 형식                      | 설명                                           |
| ------------------- | ------------------------- | ---------------------------------------------- |
| server_address      | host:port                 | ssh 서버 주소                                  |
| username            | 사용자이름                | ssh 사용자 이름                                |
| private_key_path    | 파일경로                  | ssh 사용자 개인 암호키 파일 경로               |
| remote_root         | 업로드될 경로             | 절대 경로로 입력                               |
| games               | 게임 설정                 | 동기화 대상 게임 설정                          |
| games.name          | 게임 이름                 | 중복되지 않도록 입력                           |
| games.local_dir     | 세이브 파일 폴더          | 절대 경로로 입력                               |
| games.file_patterns | 세이브 파일 필터링 정규식 | `\`문자 이스케이프에 주의해서 입력             |
| games.program_name  | 실행파일명                | (선택사항) 절대 경로 또는 파일명(filename.ext) |
