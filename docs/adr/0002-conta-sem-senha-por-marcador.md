# 0002 — Conta sem senha representada por marcador

- Status: aceito
- Data: 2026-11-13

## Contexto

Contas criadas pelo Login social não têm senha, mas `Users.password_hash` é `NOT NULL` no SQLite. O SQLite não permite remover `NOT NULL` com `ALTER TABLE`; seria preciso reconstruir a tabela (criar, copiar, renomear) na inicialização, com `PRAGMA foreign_keys = ON` ativo e `PasswordResetCodes` referenciando `Users`, e sem ferramenta de migração.

## Decisão

Representar a Conta sem senha com `password_hash = ''`. O significado do marcador fica encapsulado no domínio e no repositório (por exemplo, `HasPassword()`); nenhum outro código compara com `''`. Como o bcrypt nunca aceita um hash vazio, o login com senha numa Conta sem senha falha com "Invalid credentials".

## Consequências

- Nenhuma migração de schema arriscada na inicialização.
- O banco deixa de garantir que toda conta tem um hash real; a garantia passa para o código.
- Quando houver ferramenta de migração, o marcador pode virar `NULL` sem mudar o comportamento.
