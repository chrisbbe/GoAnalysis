package org.sonarsource.plugins.go.GoAnalyzer;

import org.sonar.api.utils.log.Logger;
import org.sonar.api.utils.log.Loggers;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.Scanner;
import java.util.concurrent.BrokenBarrierException;
import java.util.concurrent.CyclicBarrier;

class ProcessExec {
  private static final Logger LOGGER = Loggers.get(ProcessExec.class);

  private String output;
  private String errorOutput;

  void executeProcess(String[] command, String workingDir) throws IOException, InterruptedException, BrokenBarrierException {
    ProcessBuilder processBuilder = new ProcessBuilder(command);

    processBuilder.directory(new File(workingDir));
    Process process = processBuilder.start();

    CyclicBarrier cyclicBarrier = new CyclicBarrier(3); // Barrier to wait for IO to be consumed.

    IOThreadHandler outputHandler = new IOThreadHandler(process.getInputStream(), cyclicBarrier);
    IOThreadHandler errorHandler = new IOThreadHandler(process.getErrorStream(), cyclicBarrier);
    outputHandler.start();
    errorHandler.start();

    LOGGER.info("Before process.waitFor()");
    LOGGER.info("GoAnalyzer performs inspection. Please wait...");
    process.waitFor();
    LOGGER.info("After process.waitFor()");

    LOGGER.info("Before cyclicBarrier.await()");
    cyclicBarrier.await(); // Output and errorOutput should be filled!
    LOGGER.info("After cyclicBarrier.await()");

    LOGGER.info("Before output");
    this.output = outputHandler.getOutput();
    this.errorOutput = errorHandler.getOutput();
    LOGGER.info("After output");
  }

  String getOutput() {
    return output;
  }

  public String getErrorOutput() {
    return errorOutput;
  }

  private class IOThreadHandler extends Thread {
    private final InputStream inputStream;
    private final CyclicBarrier cyclicBarrier;
    private String output;

    IOThreadHandler(InputStream inputStream, CyclicBarrier cyclicBarrier) {
      this.inputStream = inputStream;
      this.cyclicBarrier = cyclicBarrier;
    }

    @Override
    public void run() {
      StringBuilder stringBuilder = new StringBuilder();

      try (Scanner br = new Scanner(new InputStreamReader(inputStream))) {
        String line;
        while (br.hasNextLine()) {
          line = br.nextLine();
          stringBuilder.append(line);
        }
      } finally {
        try {
          output = stringBuilder.toString();
          LOGGER.info("Thread " + this.getId() + " done!");
          this.cyclicBarrier.await();
        } catch (Exception e) {
          LOGGER.info(e.toString());
        }
      }

    }

    String getOutput() {
      return output;
    }
  }
}
