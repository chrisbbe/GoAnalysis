package org.sonarsource.plugins.go.rules;

import org.sonar.api.batch.fs.FileSystem;
import org.sonar.api.batch.fs.InputFile;
import org.sonar.api.batch.sensor.Sensor;
import org.sonar.api.batch.sensor.SensorContext;
import org.sonar.api.batch.sensor.SensorDescriptor;
import org.sonar.api.component.ResourcePerspectives;
import org.sonar.api.issue.Issuable;
import org.sonar.api.issue.Issue;
import org.sonar.api.measures.CoreMetrics;
import org.sonar.api.measures.Metric;
import org.sonar.api.rule.RuleKey;
import org.sonar.api.utils.log.Logger;
import org.sonar.api.utils.log.Loggers;
import org.sonarsource.plugins.go.GoAnalyzer.ExecuteGoAnalyzer;
import org.sonarsource.plugins.go.GoAnalyzer.GoFileViolation;
import org.sonarsource.plugins.go.GoAnalyzer.Violation;
import org.sonarsource.plugins.go.languages.GoLanguage;

import java.io.Serializable;
import java.util.HashMap;
import java.util.List;

import static java.lang.String.format;

/**
 * The goal of this Sensor is to load the results of an analysis performed by a fictive external tool named: FooLint
 * Results are provided as an xml file and are corresponding to the rules defined in 'rules.xml'.
 * To be very abstract, these rules are applied on source files made with the fictive language Foo.
 */
public class GoIssuesLoaderSensor implements Sensor {
  private static final Logger LOGGER = Loggers.get(GoIssuesLoaderSensor.class);
  private static final String REPORT_PATH_KEY = "sonar.go.reportPath";
  private final FileSystem fileSystem;
  private final ResourcePerspectives perspectives;
  private SensorContext sensorContext;

  /**
   * Use of IoC to get Settings, FileSystem, RuleFinder and ResourcePerspectives
   */
  public GoIssuesLoaderSensor(final FileSystem fileSystem, final ResourcePerspectives perspectives) {
    this.fileSystem = fileSystem;
    this.perspectives = perspectives;
  }

  protected String reportPathKey() {
    return REPORT_PATH_KEY;
  }

  @SuppressWarnings("PMD.ConfusingTernary")
  private boolean saveIssue(InputFile inputFile, int line, String externalRuleKey, String message) {
    RuleKey rule = RuleKey.of(GoRulesDefinition.getRepositoryKeyForLanguage(inputFile.language()), externalRuleKey);

    Issuable issuable = perspectives.as(Issuable.class, inputFile);
    boolean result = false;
    if (issuable != null) {
      LOGGER.debug("Issuable is not null: %s", issuable.toString());
      Issuable.IssueBuilder issueBuilder = issuable.newIssueBuilder()
          .ruleKey(rule)
          .message(message);
      if (line > 0) {
        LOGGER.debug("line is > 0");
        issueBuilder = issueBuilder.line(line);
      }
      Issue issue = issueBuilder.build();
      LOGGER.debug("issue == null? " + (issue == null));
      try {
        result = issuable.addIssue(issue);
        LOGGER.debug("after addIssue: result={}", result);
      } catch (org.sonar.api.utils.MessageException me) {
        LOGGER.error(format("Can't add issue on file %s at line %d.", inputFile.absolutePath(), line), me);
      }

    } else {
      LOGGER.debug("Can't find an Issuable corresponding to InputFile:" + inputFile.absolutePath());
    }
    return result;
  }

  @Override
  public void describe(final SensorDescriptor descriptor) {
    descriptor.name("GoLang Issues Loader Sensor");
    descriptor.onlyOnLanguage(GoLanguage.KEY);
  }

  @Override
  public void execute(final SensorContext sensorContext) {
    this.sensorContext = sensorContext;
    ExecuteGoAnalyzer executeGoAnalyzer = new ExecuteGoAnalyzer();
    Iterable<InputFile> goFiles = sensorContext.fileSystem().inputFiles(fileSystem.predicates().hasLanguage("go"));

    HashMap<String, InputFile> hashMap = new HashMap<>();
    for (InputFile goFile : goFiles) {
      hashMap.put(goFile.path().toString(), goFile);
    }

    LOGGER.info("Started external GoAnalyzer... :)");
    LOGGER.info("BaseDir: " + fileSystem.baseDir().getPath());
    // Execute the external tool.
    List<GoFileViolation> goFileViolationList = executeGoAnalyzer.runAnalyzer(fileSystem.baseDir().getPath());
    LOGGER.info("Got number of files: " + goFileViolationList.size());

    for (GoFileViolation goFileViolation : goFileViolationList) {
      InputFile goInputFile = hashMap.get(goFileViolation.getFilePath());
      if (goInputFile != null) {

        // Save measure metrics.
        saveMetricOnFile(goInputFile, CoreMetrics.NCLOC, goFileViolation.getLinesOfCode()); // Lines of code.

        // Save issues.
        for (Violation violation : goFileViolation.getViolations()) {
          boolean result = saveIssue(
              goInputFile,
              violation.getSrcLine(),
              violation.getType(),
              violation.getDescription()
          );
          LOGGER.debug("Issues: " + violation.getType() + " - Saved: " + result);
        }
      }
    }
    LOGGER.info("Number of Go files processed: " + hashMap.size());
  }

  private <T extends Serializable> void saveMetricOnFile(InputFile inputFile, Metric metric, T value) {
    sensorContext.<T>newMeasure()
        .withValue(value)
        .forMetric(metric)
        .on(inputFile)
        .save();
  }

}
