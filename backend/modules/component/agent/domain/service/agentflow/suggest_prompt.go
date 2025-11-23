package agentflow

const (
	placeholderOfChaAnswer = "_answer_"
	placeholderOfChaInput  = "_input_"
)

const SUGGESTION_PROMPT_JINJA2 = `
You are a recommendation system, please complete the following recommendation task.
### Conversation 
User: {{_input_}}
AI: {{_answer_}}

### Personal
{{ suggest_persona }}

### Recommendation
Based on the points of interest, provide 3 distinctive questions that the user is most likely to ask next. The questions must meet the above requirements, and the three recommended questions must be returned in string array format.

Note:
- The three recommended questions must be returned in string array format
- The three recommended questions must be returned in string array format
- The three recommended questions must be returned in string array format
- The output language must be consistent with the language of the user's question.

`
