# 📌 Contexto Geral do Projeto (para uso por IA / Cursor)

## 1. Visão do Produto

Este projeto é um **aplicativo móvel de gerenciamento de treinos de musculação**, desenvolvido em **React Native**, com backend em **ASP.NET (Web API)**.

O aplicativo permite que o usuário:

* Crie treinos para cada dia da semana
* Defina exercícios, séries, repetições, carga e intervalo
* Registre execuções reais dos treinos
* Insira resultados de avaliações de **bioimpedância**

O grande diferencial do produto é a **utilização de Inteligência Artificial por meio de Agentes de IA**, oferecidos como **recursos premium via assinatura**, com foco em personalização, acompanhamento inteligente e alto valor percebido.

---

## 2. Público-Alvo

* Praticantes de musculação iniciantes, intermediários e avançados
* Usuários que treinam sem personal trainer
* Pessoas que valorizam dados, evolução e acompanhamento inteligente

---

## 3. Objetivos da IA no Produto

A IA NÃO deve ser apenas um chatbot genérico.
Ela deve atuar como um **Coach Virtual Inteligente**, tomando decisões e fornecendo insights baseados em dados reais do usuário.

A IA deve:

* Analisar dados históricos
* Aprender padrões do usuário
* Adaptar recomendações
* Explicar decisões de forma clara e humana

---

## 4. Agentes de IA (Arquitetura Conceitual)

### 🧠 4.1 Agente Coach Virtual

Responsável por interação conversacional.

Funções:

* Responder dúvidas sobre treino
* Sugerir ajustes
* Explicar decisões da IA
* Lembrar preferências do usuário

Exemplos:

* "Posso trocar esse exercício hoje?"
* "Meu treino está pesado demais"

---

### 📊 4.2 Agente Analista de Progresso

Responsável por análise de dados históricos.

Analisa:

* Evolução de cargas
* Volume semanal
* Frequência
* Bioimpedância

Entrega:

* Relatórios em linguagem natural
* Resumos semanais/mensais

---

### 🧪 4.3 Agente de Bioimpedância

Especializado na interpretação dos dados corporais.

Dados analisados:

* Peso corporal
* Massa magra
* Gordura corporal
* Taxa metabólica basal (se disponível)

Objetivo:

* Explicar mudanças
* Relacionar com treino
* Sugerir ajustes

---

### ⚠️ 4.4 Agente de Prevenção de Lesões

Analisa riscos com base em padrões de treino.

Detecta:

* Aumento abrupto de carga
* Volume excessivo
* Pouco descanso
* Queda de performance

Entrega:

* Alertas preventivos

⚠️ Importante: nunca dar diagnóstico médico.

---

### 🗓️ 4.5 Agente de Planejamento de Treinos

Responsável por organização de ciclos.

Cria:

* Planejamento semanal
* Mesociclos (4–8 semanas)
* Deloads

Baseado em:

* Objetivo do usuário
* Disponibilidade semanal
* Nível de treino

---

## 5. Funcionalidades de IA (Premium)

### 🔥 Treinos Adaptativos Automáticos

* Ajuste de carga, reps e volume
* Baseado no desempenho real
* Considera falhas, RPE percebido e bioimpedância

---

### 💬 Chat com Memória (Coach Virtual)

* Conversa natural
* Memória curta e longa
* Contexto histórico do usuário

---

### 📈 Análise Inteligente de Progresso

* Correlação entre treino e bioimpedância
* Comparações temporais
* Explicações em texto

---

### ⚠️ Alertas de Overtraining

* Identificação de padrões perigosos
* Sugestões de ajuste

---

### 🔄 Substituição Inteligente de Exercícios

* Sugestão de equivalentes
* Mantém estímulo muscular
* Considera equipamentos disponíveis

---

### 🎯 Previsão de Resultados (Avançado)

* Estimativas baseadas em dados históricos
* Sempre apresentadas como aproximações

---

## 6. Regras Importantes para a IA

### ❌ O que a IA NÃO deve fazer

* Diagnóstico médico
* Prescrição clínica
* Garantir resultados
* Incentivar comportamentos perigosos

---

### ✅ Boas práticas

* Sempre explicar o *porquê* das sugestões
* Usar linguagem clara e motivadora
* Ser conservadora em ajustes de carga
* Usar disclaimers quando necessário

---

## 7. Estrutura Técnica (Resumo)

### Frontend

* React Native
* Telas de treino, histórico, bioimpedância e chat

### Backend

* ASP.NET Web API
* Autenticação
* Persistência de dados

### Dados importantes para IA

* Usuário
* Treinos planejados
* Execuções reais
* Bioimpedância
* Feedback subjetivo (leve / ok / pesado)

---

## 8. Modelo de Negócio

### Free

* Criação manual de treinos
* Registro de cargas
* Bioimpedância manual

### Premium (IA)

* Todos os agentes de IA
* Treinos adaptativos
* Análises inteligentes
* Planejamento automático

---

## 9. Tom de Comunicação da IA

* Profissional, motivador e acessível
* Estilo "coach experiente"
* Evitar jargões excessivos
* Foco em progresso sustentável

---

## 10. Objetivo Final

Criar um aplicativo que entregue **valor real**, fazendo o usuário sentir que tem um **personal trainer inteligente no bolso**, justificando claramente o custo da assinatura.

