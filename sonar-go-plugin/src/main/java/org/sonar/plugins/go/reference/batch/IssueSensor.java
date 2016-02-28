/*
 * SonarQube Go :: Plugin
 * Copyright (C) 2016-2016 SonarSource SA
 * mailto:contact AT sonarsource DOT com
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
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, write to the Free Software Foundation,
 * Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */
package org.sonar.plugins.go.reference.batch;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.sonar.api.batch.Sensor;
import org.sonar.api.batch.SensorContext;
import org.sonar.api.batch.fs.FileSystem;
import org.sonar.api.batch.fs.InputFile;
import org.sonar.api.component.ResourcePerspectives;
import org.sonar.api.issue.Issuable;
import org.sonar.api.resources.Project;
import org.sonar.api.rule.RuleKey;

public class IssueSensor implements Sensor {
  private static final Logger LOGGER = LoggerFactory.getLogger(IssueSensor.class);

  private final FileSystem fs;
  private final ResourcePerspectives perspectives;

  /**
   * Use of IoC to get FileSystem
   */
  public IssueSensor(FileSystem fs, ResourcePerspectives perspectives) {
    this.fs = fs;
    this.perspectives = perspectives;
  }

  @Override
  public boolean shouldExecuteOnProject(Project project) {
    // This sensor is executed only when there are Java files
    return fs.hasFiles(fs.predicates().hasLanguage("go"));
  }

  @Override
  public void analyse(Project project, SensorContext sensorContext) {
    // This sensor create an issue on each java file
    for (InputFile inputFile : fs.inputFiles(fs.predicates().hasLanguage("go"))) {
      Issuable issuable = perspectives.as(Issuable.class, inputFile);

      LOGGER.info("Added: " + issuable.addIssue(issuable.newIssueBuilder()
        .ruleKey(RuleKey.of("go", "CC"))
        .message("High cyclomatic complexity in function or method.")
        .line(1)
        .build()));
    }
  }

  @Override
  public String toString() {
    return getClass().getSimpleName();
  }

}
