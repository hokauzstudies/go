# Docker Compose

Para levantar a instância do db, basta executar:

```sh
$ docker-compose up
```

Isto irá levantar uma instância do MySQL na porta `3306` da máquina local, e uma instância do `Adminer` na
porta `8080`.

# MySQL

### Senha
Caso seja necessário, a senha para acessar a instância do mysql do docker é `admin` (definida dentro do `docker-compose.yml`)

### Populando o banco

O container automaticamente importa o arquivo de dump que contém o schema do banco e popula as tabelas com
algumas entidades:
* um super admin na tabela `users`
* um sistema na tabela `systems`
* uma operadora na tabela `operators`

Uma vez criado, o schema e os dados permanecem na pasta `/var/lib/mysql` do container

Caso seja necessário "resetar" os dados do banco, para que o schema seja importado novamente na proxima vez
que a instância for levantada com `docker-compose up`, você deve primeiro limpar a pasta de dados do
container. Caso isto não seja feita, o mysql detecta que já existem dados para o banco atual e ignora os
arquivos sql de dump.


Para limpar a pasta de dados do container, execute:

```sh
$ docker-compose run mysql /bin/bash
```

Com isso você acessará o container rodando o bash.

Depois, execute este comando no bash dentro do container:
```sh
$ rm -rf /var/lib/mysql
```

É provável que o container exiba uma mensagem dizendo que não foi possível remover a pasta pois o alvo está
ocupado, mas isto é comum.
Depois de deletar a pasta:
1) Saia do prompt atual do container, com `exit`
2) Pare os containers que estão rodando no momento com `Control+C`
3) Suba novamente a instância do banco com `docker-compose up`

Desta vez os logs avisarão que o arquivo de dump foi encontrado e irá importar o script sql que irá modelar e
popular o bando.



# Adminer

Para acessar o adminer, depois de subir os containers, basta visitar `localhost:8080` no seu navegador.
