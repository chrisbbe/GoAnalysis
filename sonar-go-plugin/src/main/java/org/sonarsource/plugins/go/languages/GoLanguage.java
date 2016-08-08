package org.sonarsource.plugins.go.languages;

import org.apache.commons.lang.StringUtils;
import org.sonar.api.config.Settings;
import org.sonar.api.resources.AbstractLanguage;

import java.util.ArrayList;
import java.util.List;

/**
 * This class defines the fictive Go language.
 */
public final class GoLanguage extends AbstractLanguage {
  public static final String KEY = "go";
  private static final String NAME = "Go";
  private static final String FILE_SUFFIXES_PROPERTY_KEY = "sonar.go.file.suffixes";
  private static final String DEFAULT_FILE_SUFFIXES = "go";

  private final Settings settings;

  public GoLanguage(Settings settings) {
    super(KEY, NAME);
    this.settings = settings;
  }

  @Override
  public String[] getFileSuffixes() {
    String[] suffixes = filterEmptyStrings(settings.getStringArray(FILE_SUFFIXES_PROPERTY_KEY));
    if (suffixes.length == 0) {
      suffixes = StringUtils.split(DEFAULT_FILE_SUFFIXES, ",");
    }
    return suffixes;
  }

  private String[] filterEmptyStrings(String[] stringArray) {
    List<String> nonEmptyStrings = new ArrayList<>();
    for (String string : stringArray) {
      if (StringUtils.isNotBlank(string.trim())) {
        nonEmptyStrings.add(string.trim());
      }
    }
    return nonEmptyStrings.toArray(new String[nonEmptyStrings.size()]);
  }

}
