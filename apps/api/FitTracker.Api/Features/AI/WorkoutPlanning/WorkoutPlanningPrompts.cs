namespace FitTracker.Api.Features.AI.WorkoutPlanning;

/// <summary>
/// Semantic Kernel prompt templates for workout planning
/// </summary>
public static class WorkoutPlanningPrompts
{
    /// <summary>
    /// System prompt that defines the AI's role and behavior
    /// </summary>
    public const string SystemPrompt = """
        You are an experienced fitness coach and workout planning specialist. Your role is to create personalized, effective workout plans based on user data and goals.

        IMPORTANT RULES:
        1. Always create safe, progressive workouts appropriate for the user's experience level
        2. Never prescribe dangerous exercises or excessive weights for beginners
        3. Consider the user's available equipment and time constraints
        4. Explain your reasoning in a motivating, professional tone
        5. Focus on compound movements for efficiency
        6. Include proper rest between muscle groups (48-72 hours)
        7. Balance push/pull movements to prevent imbalances

        You must respond ONLY with valid JSON in the exact format specified. Do not include any text before or after the JSON.
        """;

    /// <summary>
    /// Prompt template for generating a weekly workout plan
    /// </summary>
    public const string GenerateWeeklyPlanPrompt = """
        Based on the following user profile, create a personalized weekly workout plan:

        USER PROFILE:
        - Name: {{$userName}}
        - Age: {{$userAge}}
        - Current Weight: {{$currentWeight}} kg
        - Goal Weight: {{$goalWeight}} kg
        - Experience Level: {{$experienceLevel}}
        - Available Days per Week: {{$availableDays}}
        - Preferred Workout Duration: {{$preferredDuration}} minutes
        - Available Equipment: {{$equipment}}
        - Primary Goal: {{$goal}}

        HISTORICAL DATA:
        - Total Completed Workouts: {{$totalWorkouts}}
        - Workouts This Month: {{$monthlyWorkouts}}
        - Latest Body Fat: {{$bodyFat}}%
        - Latest Muscle Mass: {{$muscleMass}} kg
        - Current Exercises They Do: {{$currentExercises}}

        ADDITIONAL NOTES FROM USER:
        {{$additionalNotes}}

        Generate a complete weekly workout plan with the following JSON structure:
        {
            "summary": "Brief 1-2 sentence overview of the plan",
            "rationale": "2-3 sentences explaining why this plan is suitable for the user",
            "workouts": [
                {
                    "name": "Workout name (e.g., 'Push Day A')",
                    "description": "Brief description of the workout focus",
                    "goal": "Hypertrophy|Strength|Endurance|WeightLoss|Maintenance",
                    "days": ["Monday", "Thursday"],
                    "estimatedDurationMinutes": 45,
                    "exercises": [
                        {
                            "name": "Exercise name",
                            "muscleGroup": "Chest|Back|Shoulders|Biceps|Triceps|Legs|Glutes|Abs|Cardio",
                            "sets": 3,
                            "reps": 12,
                            "suggestedWeight": 20.0,
                            "restSeconds": 90,
                            "notes": "Optional technique notes"
                        }
                    ]
                }
            ]
        }

        Remember:
        - For beginners: 2-3 exercises per muscle group, moderate weights, focus on form
        - For intermediate: 3-4 exercises per muscle group, progressive overload
        - For advanced: 4-5 exercises per muscle group, varied rep ranges, intensity techniques
        - If equipment is "Minimal", use only bodyweight exercises
        - If equipment is "Home", use dumbbells, resistance bands, and bodyweight
        - If equipment is "Gym", use full range of equipment

        Respond with ONLY the JSON, no additional text.
        """;

    /// <summary>
    /// Prompt for explaining a generated plan
    /// </summary>
    public const string ExplainPlanPrompt = """
        You created the following workout plan for a user:

        {{$planJson}}

        The user wants to understand why you made these choices. Provide a detailed but accessible explanation covering:
        1. Why you chose these specific exercises
        2. Why the sets/reps scheme matches their goal
        3. How the weekly split optimizes recovery
        4. Tips for getting the most out of this plan

        Keep the tone motivating and professional, like an experienced coach talking to their client.
        """;

    /// <summary>
    /// Helper to convert enums to user-friendly strings for prompts
    /// </summary>
    public static string GetGoalDescription(Workouts.WorkoutGoal goal) => goal switch
    {
        Workouts.WorkoutGoal.Hypertrophy => "Build muscle mass and size",
        Workouts.WorkoutGoal.Strength => "Increase strength and power",
        Workouts.WorkoutGoal.Endurance => "Improve muscular endurance and stamina",
        Workouts.WorkoutGoal.WeightLoss => "Lose body fat while preserving muscle",
        Workouts.WorkoutGoal.Maintenance => "Maintain current fitness level",
        _ => "General fitness improvement"
    };

    public static string GetExperienceDescription(Profiles.ExperienceLevel level) => level switch
    {
        Profiles.ExperienceLevel.Beginner => "Beginner (less than 1 year of consistent training)",
        Profiles.ExperienceLevel.Intermediate => "Intermediate (1-3 years of consistent training)",
        Profiles.ExperienceLevel.Advanced => "Advanced (3+ years of consistent training)",
        _ => "Unknown experience level"
    };

    public static string GetEquipmentDescription(Profiles.EquipmentType equipment) => equipment switch
    {
        Profiles.EquipmentType.Gym => "Full gym (barbells, dumbbells, cables, machines)",
        Profiles.EquipmentType.Home => "Home setup (dumbbells, resistance bands, pull-up bar)",
        Profiles.EquipmentType.Minimal => "Minimal/bodyweight only",
        _ => "Unknown equipment"
    };
}
