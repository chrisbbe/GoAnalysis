/*
 * SonarQube Protocol Buffers Plugin
 * Copyright (C) 2015 SonarSource
 * sonarqube@googlegroups.com
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 3 of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02
 */
package org.sonarsource.plugins.go.languages;

import org.sonar.api.profiles.ProfileDefinition;
import org.sonar.api.profiles.RulesProfile;
import org.sonar.api.rules.Rule;
import org.sonar.api.rules.RuleFinder;
import org.sonar.api.rules.RuleQuery;
import org.sonar.api.utils.ValidationMessages;
import org.sonarsource.plugins.go.rules.GoRulesDefinition;

/**
 * Default Quality profile for the projects having files of language "go"
 */
public final class GoQualityProfile extends ProfileDefinition {
  private static final String REPOSITORY_KEY = GoRulesDefinition.getRepositoryKeyForLanguage(GoLanguage.KEY);

  private final RuleFinder ruleFinder;

  public GoQualityProfile(RuleFinder ruleFinder) {
    this.ruleFinder = ruleFinder;
  }

  @Override
  public RulesProfile createProfile(ValidationMessages validation) {
    RulesProfile ruleProfile = RulesProfile.create("Go Rules", GoLanguage.KEY);

    for (Rule rule : ruleFinder.findAll(RuleQuery.create().withRepositoryKey(REPOSITORY_KEY))) {
      ruleProfile.activateRule(rule, null);
    }
    return ruleProfile;
  }
}
