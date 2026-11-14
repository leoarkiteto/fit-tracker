# 0001 — Vinculação automática pelo E-mail verificado do provedor

- Status: aceito
- Data: 2026-11-13

## Contexto

O Login social (Google e GitHub) pode chegar com um e-mail que já pertence a uma conta do FitTracker criada com senha. Havia três caminhos: recusar o login, pedir a senha do FitTracker antes de vincular, ou vincular automaticamente.

## Decisão

Vincular automaticamente a Identidade à conta existente **somente** quando o provedor afirma que o e-mail é verificado (Google: `email_verified`; GitHub: e-mail primário e verificado de `/user/emails`). Sem E-mail verificado, o login é recusado. A comparação de e-mail ignora maiúsculas e minúsculas; se houver mais de uma conta com o mesmo e-mail, nada é vinculado e o login é recusado. Depois da primeira vinculação, a pessoa é reconhecida só pelo provedor + identificador; o e-mail do provedor não é sincronizado.

Quando a vinculação automática liga a Identidade a uma conta existente **que tem senha**, a senha é removida no mesmo passo (a conta passa a ser uma Conta sem senha, ver [ADR 0002](0002-conta-sem-senha-por-marcador.md)) e os Códigos de redefinição pendentes dessa conta são invalidados. Isso acontece só na vinculação automática, não nos logins seguintes pela Identidade já vinculada.

## Consequências

- Quem controla uma conta Google/GitHub com o e-mail verificado ganha acesso à conta do FitTracker com aquele e-mail. Aceitamos isso porque o mesmo controle do e-mail já basta para a Recuperação de senha.
- Não há passo extra de confirmação, o que simplifica a experiência e o fluxo.
- Trocar para "pedir a senha antes de vincular" depois exigiria lidar com vínculos já criados sem essa prova.
- **Pre-account hijacking:** o cadastro com senha não verifica o e-mail, então um atacante pode cadastrar o e-mail da vítima com uma senha dele antes de ela usar o Login social. Sem defesa, a vinculação automática ligaria a Identidade da vítima à conta do atacante, que continuaria entrando com a senha dele. Remover a senha na vinculação automática fecha esse caminho: a senha não provou a posse do e-mail, o provedor provou.
- O dono legítimo que já tinha senha perde o login com senha na primeira vez que usa o Login social; ele pode definir a senha de novo pela Recuperação de senha, que também prova a posse do e-mail.
- Invalidar os Códigos de redefinição pendentes impede que um código pedido antes da vinculação seja usado depois dela.
- Risco residual: sessões (JWT) já emitidas para a conta continuam válidas até expirar (24 h), porque o FitTracker não revoga tokens.
- Dados criados na conta antes da vinculação (perfil, treinos) continuam nela.
