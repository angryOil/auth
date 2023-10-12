# 간단한 auth서버 입니다.

간단한 인증,인가 ,회원가입기능을 합니다 . 


```shell
# local postgres run (docker-compose)
make local-db
# local postgres migrate init
make local-init
# local postgres apply migrate
make local-migrate
```