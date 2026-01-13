# 📌 Project General Context (for AI / Cursor usage)

## 1. Product Vision

This project is a **mobile workout management application**, developed in **React Native**, with backend in **ASP.NET (Web API)**.

The application allows the user to:

* Create workouts for each day of the week
* Define exercises, sets, reps, weight, and rest intervals
* Track actual workout executions
* Record **bioimpedance** assessment results

The product differentiator is the **use of Artificial Intelligence through AI Agents**, offered as **premium features via subscription**, focusing on personalization, intelligent monitoring, and high perceived value.

---

## 2. Target Audience

* Beginner, intermediate, and advanced weight training practitioners
* Users who train without a personal trainer
* People who value data, progress, and intelligent monitoring

---

## 3. AI Objectives in the Product

AI should NOT be just a generic chatbot.
It should act as an **Intelligent Virtual Coach**, making decisions and providing insights based on real user data.

The AI should:

* Analyze historical data
* Learn user patterns
* Adapt recommendations
* Explain decisions clearly and in a human way

---

## 4. AI Agents (Conceptual Architecture)

### 🧠 4.1 Virtual Coach Agent

Responsible for conversational interaction.

Functions:

* Answer training questions
* Suggest adjustments
* Explain AI decisions
* Remember user preferences

Examples:

* "Can I switch this exercise today?"
* "My workout is too heavy"

---

### 📊 4.2 Progress Analyst Agent

Responsible for historical data analysis.

Analyzes:

* Weight progression
* Weekly volume
* Frequency
* Bioimpedance

Delivers:

* Reports in natural language
* Weekly/monthly summaries

---

### 🧪 4.3 Bioimpedance Agent

Specialized in body data interpretation.

Data analyzed:

* Body weight
* Lean mass
* Body fat
* Basal metabolic rate (if available)

Objective:

* Explain changes
* Relate to training
* Suggest adjustments

---

### ⚠️ 4.4 Injury Prevention Agent

Analyzes risks based on training patterns.

Detects:

* Abrupt weight increase
* Excessive volume
* Insufficient rest
* Performance drop

Delivers:

* Preventive alerts

⚠️ Important: never give medical diagnosis.

---

### 🗓️ 4.5 Workout Planning Agent

Responsible for cycle organization.

Creates:

* Weekly planning
* Mesocycles (4-8 weeks)
* Deloads

Based on:

* User goal
* Weekly availability
* Training level

---

## 5. AI Features (Premium)

### 🔥 Automatic Adaptive Workouts

* Weight, reps, and volume adjustment
* Based on actual performance
* Considers failures, perceived RPE, and bioimpedance

---

### 💬 Chat with Memory (Virtual Coach)

* Natural conversation
* Short and long-term memory
* User historical context

---

### 📈 Smart Progress Analysis

* Workout and bioimpedance correlation
* Temporal comparisons
* Text explanations

---

### ⚠️ Overtraining Alerts

* Dangerous pattern identification
* Adjustment suggestions

---

### 🔄 Smart Exercise Substitution

* Equivalent suggestions
* Maintains muscle stimulus
* Considers available equipment

---

### 🎯 Results Prediction (Advanced)

* Estimates based on historical data
* Always presented as approximations

---

## 6. Important Rules for AI

### ❌ What AI should NOT do

* Medical diagnosis
* Clinical prescription
* Guarantee results
* Encourage dangerous behaviors

---

### ✅ Best practices

* Always explain the *why* of suggestions
* Use clear and motivating language
* Be conservative in weight adjustments
* Use disclaimers when necessary

---

## 7. Technical Structure (Summary)

### Frontend

* React Native
* Workout, history, bioimpedance, and chat screens

### Backend

* ASP.NET Web API
* Authentication
* Data persistence

### Important data for AI

* User
* Planned workouts
* Actual executions
* Bioimpedance
* Subjective feedback (light / ok / heavy)

---

## 8. Business Model

### Free

* Manual workout creation
* Weight tracking
* Manual bioimpedance

### Premium (AI)

* All AI agents
* Adaptive workouts
* Smart analytics
* Automatic planning

---

## 9. AI Communication Tone

* Professional, motivating, and accessible
* "Experienced coach" style
* Avoid excessive jargon
* Focus on sustainable progress

---

## 10. Final Goal

Create an application that delivers **real value**, making the user feel they have a **smart personal trainer in their pocket**, clearly justifying the subscription cost.
