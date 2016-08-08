package org.sonarsource.plugins.go.GoAnalyzer;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import org.sonar.api.utils.log.Logger;
import org.sonar.api.utils.log.Loggers;

import java.io.File;
import java.lang.reflect.Type;
import java.net.URI;
import java.util.Arrays;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;

public class ExecuteGoAnalyzer {
  private static final Logger LOGGER = Loggers.get(ExecuteGoAnalyzer.class);

  private final HashMap<String, String> analyzers;

  public ExecuteGoAnalyzer() {
    this.analyzers = new HashMap<>();
    this.analyzers.put("Linux", "analyzer/GoAnalyzerLinux");
    this.analyzers.put("Mac", "analyzer/GoAnalyzerMac");
    this.analyzers.put("Windows", "analyzer/GoAnalyzerWindows.exe");
  }

  public List<GoFileViolation> runAnalyzer(final String filePath) {
    LOGGER.info("Starting runAnalyzer");
    JarExtractor jarExtractor = new JarExtractor();
    Gson gson = new Gson();
    LOGGER.info("Done creating jar and gson objects");

    List<GoFileViolation> goFileViolations = null;

    String operatingSystem = System.getProperty("os.name");
    String architecture = System.getProperty("os.arch");
    LOGGER.info("OS: " + operatingSystem);
    LOGGER.info("Architecture: " + architecture);

    String analyzerFile = this.analyzers.get(operatingSystem.split(" ")[0]);
    if (analyzerFile == null) {
      LOGGER.error("OS " + operatingSystem + " is not supported! Exiting!");
      System.exit(1); //TODO: Should maybe throw an exception or something, not nice!
    }

    LOGGER.info("AnalyzerFile: " + analyzerFile);
    try {
      final URI exe = jarExtractor.extractFileFromJar(analyzerFile);

      String[] command = new String[4];
      command[0] = exe.getPath();
      command[1] = "-json"; // We want JSON output.
      command[2] = "-dir";
      command[3] = filePath;

      File workingDir = new File(exe.getPath()).getParentFile();
      ProcessExec processExec = new ProcessExec();
      LOGGER.info("Before executing process");
      LOGGER.info("Command: " + Arrays.toString(command));
      processExec.executeProcess(command, workingDir.getPath());
      LOGGER.info("After executing process");

      LOGGER.info("Error output: " + processExec.getErrorOutput());
      if (processExec.getErrorOutput().length() > 0) {
        LOGGER.info("###################################### GO ANALYZER ERRORS ######################################");
        LOGGER.info(processExec.getErrorOutput());
        System.exit(1); //TODO: Should maybe throw an exception or something, not nice!
      }

      Type violationType = new TypeToken<Collection<GoFileViolation>>() {
      }.getType();
      LOGGER.info("Before converting JSON to objects");
      goFileViolations = gson.fromJson(processExec.getOutput(), violationType);
      LOGGER.info("After converting JSON to objects: Number of GoFiles: " + goFileViolations.size());

    } catch (Exception e) {
      LOGGER.info(e.toString());
    }
    return goFileViolations;
  }

}
