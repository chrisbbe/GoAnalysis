package org.sonarsource.plugins.go.rules;


class GoIssue {
  private final String type;
  private final String description;
  private final String filePath;
  private final int line;

  // String type is actual, externalRuleKey.
  public GoIssue(final String type, final String description, final String filePath, final int line) {
    this.type = type;
    this.description = description;
    this.filePath = filePath;
    this.line = line;
  }

  public String getType() {
    return type;
  }

  public String getDescription() {
    return description;
  }

  public String getFilePath() {
    return filePath;
  }

  public int getLine() {
    return line;
  }

  @Override
  public String toString() {
    return type +
        "|" +
        description +
        "|" +
        filePath +
        "(" +
        line +
        ")";
  }
}
