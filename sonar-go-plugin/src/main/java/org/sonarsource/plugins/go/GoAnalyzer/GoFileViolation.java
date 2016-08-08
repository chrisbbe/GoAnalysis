package org.sonarsource.plugins.go.GoAnalyzer;

// This class maps JSON from the external GoAnalyzer tool into Java objects,
// that's the reason for capital variables as they correspond against JSON variables.
@SuppressWarnings("PMD.VariableNamingConventions")
public class GoFileViolation {
  private String FilePath;
  private int LinesOfCode;
  private int LinesOfComments;
  private Violation[] Violations;

  public String getFilePath() {
    return FilePath;
  }

  public int getLinesOfCode() {
    return LinesOfCode;
  }

  public int getLinesOfComments() {
    return LinesOfComments;
  }

  public Violation[] getViolations() {
    return Violations;
  }

}
