package org.sonarsource.plugins.go;

import org.sonar.api.Plugin;
import org.sonarsource.plugins.go.languages.GoLanguage;
import org.sonarsource.plugins.go.languages.GoQualityProfile;
import org.sonarsource.plugins.go.rules.GoIssuesLoaderSensor;
import org.sonarsource.plugins.go.rules.GoRulesDefinition;

/**
 * This class is the entry point for all extensions. It is referenced in pom.xml.
 */
public class GoPlugin implements Plugin {

  @Override
  public void define(Context context) {
    // Language.
    context.addExtensions(GoLanguage.class, GoQualityProfile.class);

    // Rules.
    context.addExtensions(GoRulesDefinition.class, GoIssuesLoaderSensor.class);
  }
}
