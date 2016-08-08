package org.sonarsource.plugins.go.GoAnalyzer;

import org.sonarsource.plugins.go.GoPlugin;

import java.io.*;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.security.CodeSource;
import java.security.ProtectionDomain;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

class JarExtractor {

  private static URI getJarURI() throws URISyntaxException {
    final ProtectionDomain domain;
    final CodeSource source;
    final URL url;
    final URI uri;

    domain = GoPlugin.class.getProtectionDomain();
    source = domain.getCodeSource();
    url = source.getLocation();
    uri = url.toURI();

    return uri;
  }

  private static URI getFile(final URI where, final String fileName) throws IOException {
    final File location = new File(where);
    final URI fileURI;

    // not in a JAR, just return the path on disk
    if (location.isDirectory()) {
      fileURI = URI.create(where.toString() + fileName);
    } else {
      final ZipFile zipFile;

      zipFile = new ZipFile(location);

      try {
        fileURI = extract(zipFile, fileName);
      } finally {
        zipFile.close();
      }
    }

    return fileURI;
  }

  private static URI extract(final ZipFile zipFile, final String fileName) throws IOException {
    final File tempFile;
    final ZipEntry entry;
    final InputStream zipStream;
    OutputStream fileStream;

    tempFile = File.createTempFile(fileName, Long.toString(System.currentTimeMillis()));
    tempFile.deleteOnExit();
    entry = zipFile.getEntry(fileName);

    if (entry == null) {
      throw new FileNotFoundException("cannot find file: " + fileName + " in archive: " + zipFile.getName());
    }

    zipStream = zipFile.getInputStream(entry);
    fileStream = null;

    try {
      final byte[] buf;
      int i;

      fileStream = new FileOutputStream(tempFile);
      buf = new byte[1024];

      while ((i = zipStream.read(buf)) != -1) {
        fileStream.write(buf, 0, i);
      }

      tempFile.setReadable(true);
      tempFile.setExecutable(true);

    } finally {
      zipStream.close();
      if (fileStream != null) {
        fileStream.close();
      }
    }

    return tempFile.toURI();
  }

  public URI extractFileFromJar(final String fileName) throws URISyntaxException, IOException {
    return getFile(getJarURI(), fileName);
  }


}

