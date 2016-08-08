package org.sonarsource.plugins.go.rules;

import org.sonar.api.server.rule.RulesDefinition;
import org.sonar.api.server.rule.RulesDefinitionXmlLoader;
import org.sonarsource.plugins.go.languages.GoLanguage;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.Locale;

public final class GoRulesDefinition implements RulesDefinition {
  private static final String KEY = "go";
  private static final String NAME = "Go";

  public static String getRepositoryKeyForLanguage(String languageKey) {
    return languageKey.toLowerCase(Locale.ENGLISH) + "-" + KEY;
  }

  private static String getRepositoryNameForLanguage(String languageKey) {
    return languageKey.toUpperCase(Locale.ENGLISH) + " " + NAME;
  }

  private String rulesDefinitionFilePath() {
    return "/ruleset/go-rules.xml";
  }

  private void defineRulesForLanguage(Context context, String repositoryKey, String repositoryName, String languageKey) {
    NewRepository repository = context.createRepository(repositoryKey, languageKey).setName(repositoryName);

    InputStream rulesXml = this.getClass().getResourceAsStream(rulesDefinitionFilePath());
    if (rulesXml != null) {
      RulesDefinitionXmlLoader rulesLoader = new RulesDefinitionXmlLoader();
      rulesLoader.load(repository, rulesXml, StandardCharsets.UTF_8.name());
    }

    repository.done();
  }

  @Override
  public void define(Context context) {
    String repositoryKey = GoRulesDefinition.getRepositoryKeyForLanguage(GoLanguage.KEY);
    String repositoryName = GoRulesDefinition.getRepositoryNameForLanguage(GoLanguage.KEY);
    defineRulesForLanguage(context, repositoryKey, repositoryName, GoLanguage.KEY);
  }

}
