# Glossário

Linguagem do domínio do FitTracker. Cada termo traz, entre parênteses, o nome usado no código.

## Conta e acesso

**Recuperação de senha** (Password Recovery)
Processo pelo qual uma pessoa que esqueceu a senha prova que é dona do e-mail da conta e define uma nova senha, sem estar logada. É diferente da **Troca de senha**.
_Evitar_: "reset de conta", "esqueci a conta".

**Troca de senha** (Change Password)
Substituição da senha por alguém já logado, que informa a senha atual. Não faz parte da Recuperação de senha.

**Pedido de redefinição** (Password Reset Request)
Solicitação, feita a partir do e-mail, para iniciar a Recuperação de senha. A resposta é sempre a mesma, exista ou não uma conta com aquele e-mail.

**Código de redefinição** (Password Reset Code)
Código de 6 dígitos enviado ao e-mail da conta, que prova a posse do e-mail. Vale por 15 minutos, pode ser usado uma única vez e deixa de valer após 5 tentativas erradas ou quando um novo Código de redefinição é emitido para a mesma conta.
_Evitar_: "token", "senha temporária".

**Intervalo de reenvio** (Resend Cooldown)
Tempo mínimo de 60 segundos entre dois Pedidos de redefinição para o mesmo e-mail. Um pedido feito dentro do intervalo não gera um novo Código de redefinição; o anterior continua valendo.
