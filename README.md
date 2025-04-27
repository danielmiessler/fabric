# MCP-Server-Fabric

![image.png](image.png)

![image.png](image%201.png)

![image.png](image%202.png)

![image.png](image%203.png)

# Fabricer Project README

## Overview

Fabricer integrates patterns from Fabric as tools in N8N, enabling customization and enhancement of each pattern. This setup provides seamless access to various integrations, leveraging the MCP-Server to utilize tools as a USB, simplifying pattern integration and empowering AI capabilities significantly.

## Key Features

- **MCP-Server Integration**: Using MCP-Server with Fabric in Claude simplifies AI interaction with the model context protocol, utilizing diverse tools to create a precise and intelligent chain of thought output.
- **Prompt Filtering by Tags**: Filter prompts using tags such as AI, DEVELOPMENT, ANALYSIS, etc., for targeted tool application.
- **Customization and Enhancement**: Customize and improve each pattern within N8N, enhancing functionality through integrations.
- **Supercharged AI**: Integrating patterns as tools feels like giving superpowers to the AI, enhancing its capabilities.

## Setup Instructions

1. **How to start**: Import Fabricer workflow into your n8n workplace
2. **Update Credentials**: Update API credentials to ensure smooth operation.
3. **Update parameters**: As your preferences as tags or patterns about to use
4. **start workflow** the last node to avoid any inconveniences
5. **Error Handling**: If you encounter errors like "write EPROTO F8212C9BC37F0000:error:0A00010B:SSL routines:ssl3_get_record:wrong version number", restart N8N from Docker. This issue arises due to N8N API limitations, which cap at 120 calls simultaneously.

## Usage

- Use the MCP-Server to integrate and run tools, enhancing AI output precision.
- Filter prompts by tags to apply the right tools for your needs.
- Customize patterns in N8N to fit specific requirements, leveraging available integrations.

Using the patterns from Fabric as tools in N8N allow me to customize and  improve each pattern, giving access to different integrations and using it with a MCP-Server, it allow me to use all the tools as a USB making easy the integration of each pattern as a tool, feels like giving super powers to the AI

Using the MCP-Server -Fabric in Claude i can notice how easy is for the ai to work with the model context protocol using different tools, creating the chain of thought  bringing a output more precise and smart 

[https://claude.ai/share/e394a3ad-0185-42a6-9a73-b9a1b413ad2c](https://claude.ai/share/e394a3ad-0185-42a6-9a73-b9a1b413ad2c)

### List of Tags

- analysis
- brainstorming
- career
- coding
- communication
- content_creation
- creativity
- cybersecurity
- data_analysis
- decision_making
- education
- ethics
- fiction
- fun
- health
- humor
- learning
- legal
- marketing
- personal_development
- philosophy
- productivity
- research
- science
- security
- social_media
- storytelling
- strategy
- teaching
- technology
- threat_modeling
- writing

### List of pattern

agility_story
ai
analyze_answers
analyze_bill
analyze_bill_short
analyze_candidates
analyze_cfp_submission
analyze_claims
analyze_comments
analyze_debate
analyze_email_headers
analyze_incident
analyze_interviewer_techniques
analyze_logs
analyze_malware
analyze_military_strategy
analyze_mistakes
analyze_paper
analyze_patent
analyze_personality
analyze_presentation
analyze_product_feedback
analyze_proposition
analyze_prose
analyze_prose_json
analyze_prose_pinker
analyze_risk
analyze_sales_call
analyze_spiritual_text
analyze_tech_impact
analyze_threat_report
analyze_threat_report_cmds
analyze_threat_report_trends
answer_interview_question
ask_secure_by_design_questions
ask_uncle_duke
capture_thinkers_work
check_agreement
clean_text
coding_master
compare_and_contrast
convert_to_markdown
create_5_sentence_summary
create_academic_paper
create_ai_jobs_analysis
create_aphorisms
create_art_prompt
create_better_frame
create_coding_feature
create_coding_project
create_command
create_cyber_summary
create_design_document
create_diy
create_excalidraw_visualization
create_flash_cards
create_formal_email
create_git_diff_commit
create_graph_from_input
create_hormozi_offer
create_idea_compass
create_investigation_visualization
create_keynote
create_loe_document
create_logo
create_markmap_visualization
create_mermaid_visualization
create_mermaid_visualization_for_github
create_micro_summary
create_network_threat_landscape
create_newsletter_entry
create_npc
create_pattern
create_prd
create_prediction_block
create_quiz
create_reading_plan
create_recursive_outline
create_report_finding
create_rpg_summary
create_security_update
create_show_intro
create_sigma_rules
create_story_explanation
create_stride_threat_model
create_summary
create_tags
create_threat_scenarios
create_ttrc_graph
create_ttrc_narrative
create_upgrade_pack
create_user_story
create_video_chapters
create_visualization
dialog_with_socrates
enrich_blog_post
explain_code
explain_docs
explain_math
explain_project
explain_terms
export_data_as_csv
extract_algorithm_update_recommendations
extract_article_wisdom
extract_book_ideas
extract_book_recommendations
extract_business_ideas
extract_controversial_ideas
extract_core_message
extract_ctf_writeup
extract_domains
extract_extraordinary_claims
extract_ideas
extract_insights
extract_insights_dm
extract_instructions
extract_jokes
extract_latest_video
extract_main_activities
extract_main_idea
extract_most_redeeming_thing
extract_patterns
extract_poc
extract_predictions
extract_primary_problem
extract_primary_solution
extract_product_features
extract_questions
extract_recipe
extract_recommendations
extract_references
extract_skills
extract_song_meaning
extract_sponsors
extract_videoid
extract_wisdom
extract_wisdom_agents
extract_wisdom_dm
extract_wisdom_nometa
find_female_life_partner
find_hidden_message
find_logical_fallacies
get_wow_per_minute
get_youtube_rss
humanize
identify_dsrp_distinctions
identify_dsrp_perspectives
identify_dsrp_relationships
identify_dsrp_systems
identify_job_stories
improve_academic_writing
improve_prompt
improve_report_finding
improve_writing
judge_output
label_and_rate
md_callout
official_pattern_template
prepare_7s_strategy
provide_guidance
rate_ai_response
rate_ai_result
rate_content
rate_value
raw_query
raycast
recommend_artists
recommend_pipeline_upgrades
recommend_talkpanel_topics
refine_design_document
review_design
sanitize_broken_html_to_markdown
show_fabric_options_markmap
solve_with_cot
suggest_pattern
summarize
summarize_debate
summarize_git_changes
summarize_git_diff
summarize_lecture
summarize_legislation
summarize_meeting
summarize_micro
summarize_newsletter
summarize_paper
summarize_prompt
summarize_pull-requests
summarize_rpg_session
t_analyse_challenge_handling
t_check_metrics
t_create_h3_career
t_create_opening_sentences
t_describe_life_outlook
t_extract_intro_sentences
t_extract_panel_topics
t_find_blindspots
t_find_negative_thinking
t_find_neglected_goals
t_give_encouragement
t_red_team_thinking
t_threat_model_plans
t_visualize_mission_goals_projects
t_year_in_review
to_flashcards
transcribe_minutes
translate
tweet
write_essay
write_hackerone_report
write_latex
write_micro_essay
write_nuclei_template_rule
write_pull-request
write_semgrep_rule
youtube_summary
