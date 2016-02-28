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

import org.sonar.api.batch.CheckProject;
import org.sonar.api.batch.SensorContext;
import org.sonar.api.issue.Issue;
import org.sonar.api.issue.ProjectIssues;
import org.sonar.api.resources.Project;

public class ListAllIssuesPostJob implements org.sonar.api.batch.PostJob, CheckProject {

  private final ProjectIssues projectIssues;

  public ListAllIssuesPostJob(ProjectIssues projectIssues) {
    this.projectIssues = projectIssues;
  }

  @Override
  public boolean shouldExecuteOnProject(Project project) {
    return Boolean.TRUE;
  }

  @Override
  public void executeOn(Project project, SensorContext context) {
    System.out.println("ListAllIssuesPostJob");

    // all open issues
    for (Issue issue : projectIssues.issues()) {
      String ruleKey = issue.ruleKey().toString();
      Integer issueLine = issue.line();
      String severity = issue.severity();
      boolean isNew = issue.isNew();

      // just to illustrate, we dump some fields of the 'issue' in sysout (bad, very bad)
      System.out.println(ruleKey + " : " + issue.componentKey() + "(" + issueLine + ")");
      System.out.println("isNew: " + isNew + " | severity: " + severity);
    }

    // all resolved issues
    for (Issue issue : projectIssues.resolvedIssues()) {
      String ruleKey = issue.ruleKey().toString();
      Integer issueLine = issue.line();
      boolean isNew = issue.isNew();

      System.out.println(ruleKey + " : " + issue.componentKey() + "(" + issueLine + ")");
      System.out.println("isNew: " + isNew + " | resolution: " + issue.resolution());
    }
  }
}
