resource "dynatrace_user_experience_score" "#name#" {
  consider_last_action                  = false
  consider_rage_click                   = false
  max_frustrated_user_actions_threshold = 20
  min_satisfied_user_actions_threshold  = 60
}
