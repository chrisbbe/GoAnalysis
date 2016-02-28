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
package org.sonar.plugins.go;

import org.sonar.api.Properties;
import org.sonar.api.Property;
import org.sonar.api.SonarPlugin;
  import org.sonar.plugins.go.reference.batch.IssueSensor;
import org.sonar.plugins.go.reference.batch.ListAllIssuesPostJob;
import org.sonar.plugins.go.reference.rules.GoProfile;
import org.sonar.plugins.go.reference.rules.GoRulesDefinition;

import java.util.Arrays;
import java.util.List;

/**
 * This class is the entry point for all extensions
 */
@Properties({
    @Property(
        key = GoPlugin.MY_PROPERTY,
        name = "Plugin Property",
        description = "A property for the plugin",
        defaultValue = "Hello World!"),
    @Property(
        key = GoPlugin.FILE_SUFFIXES_KEY,
        name = "File Suffixes",
        description = "Comma-separated list of suffixes for files to analyze.",
        defaultValue = GoPlugin.DEFAULT_FILE_SUFFIXES)
})
public final class GoPlugin extends SonarPlugin {
  public static final String MY_PROPERTY = "sonar.example.myproperty";

  public static final String FILE_SUFFIXES_KEY = "sonar.go.file.suffixes";
  public static final String DEFAULT_FILE_SUFFIXES = "go";

  // This is where you're going to declare all your SonarQube extensions
  @Override
  public List getExtensions() {
    return Arrays.asList(
        // Language
        GoLanguage.class,

        // Rules, Quality Profile
        GoRulesDefinition.class, GoProfile.class,

        // Batch
        IssueSensor.class, ListAllIssuesPostJob.class
    );
  }
}
