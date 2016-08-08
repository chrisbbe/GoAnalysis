package org.sonarsource.plugins.go.GoAnalyzer;

// This class maps JSON from the external GoAnalyzer tool into Java objects,
// that's the reason for capital variables as they correspond against JSON variables.
@SuppressWarnings("PMD.VariableNamingConventions")
public class Violation {
  private String Type;
  private String Description;
  private int SrcLine;

  public String getType() {
    return Type;
  }

  public void setType(String type) {
    Type = type;
  }

  public String getDescription() {
    return Description;
  }

  public void setDescription(String description) {
    Description = description;
  }

  public int getSrcLine() {
    return SrcLine;
  }

  public void setSrcLine(int srcLine) {
    SrcLine = srcLine;
  }
}
