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

**Login social** (Social Sign-In)
Entrar no FitTracker com uma conta pessoal de um Provedor de identidade, em vez de e-mail e senha. Na primeira vez, cria a conta ou a vincula a uma conta existente com o mesmo E-mail verificado; nas seguintes, a pessoa é reconhecida pela Identidade vinculada, e o e-mail do provedor não altera a conta.
_Evitar_: "SSO" (reservado para login único corporativo, que o FitTracker não oferece).

**Provedor de identidade** (Identity Provider)
Serviço externo que atesta quem a pessoa é, para que ela entre no FitTracker sem digitar a senha do FitTracker. Hoje: Google e GitHub. O provedor só prova a identidade; a sessão continua sendo do FitTracker.

**Identidade vinculada** (Linked Identity)
Associação permanente entre uma conta do FitTracker e uma conta em um Provedor de identidade, identificada pelo provedor e pelo identificador que ele dá à pessoa. Cada conta de provedor pertence a no máximo uma conta do FitTracker; uma conta do FitTracker pode ter várias Identidades vinculadas.

**E-mail verificado** (Verified Email)
E-mail que o Provedor de identidade afirma pertencer à pessoa. Só um E-mail verificado permite criar uma conta ou vincular uma Identidade a uma conta existente.
_Evitar_: "e-mail confirmado" quando a confirmação foi feita pelo FitTracker.

**Conta sem senha** (Password-less Account)
Conta sem senha do FitTracker: criada pelo Login social, ou que perdeu a senha ao ser vinculada automaticamente a uma Identidade. Só entra pelo provedor, até definir uma senha pela Recuperação de senha. A Troca de senha não se aplica a ela.
