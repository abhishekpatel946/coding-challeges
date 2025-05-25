using System;
using System.IO;
using System.Text;

class Program
{
    static void Main(string[] args)
    {
        bool useStdin = args.Length == 1 && args[0].StartsWith("-");
        string option = useStdin ? args[0] : (args.Length > 0 ? args[0] : "");
        string filePath = useStdin ? null : (args.Length == 2 ? args[1] : args.Length == 1 ? args[0] : null);

        if (filePath == null && !useStdin)
        {
            Console.WriteLine("Usage: ccwc [-c|-l|-w|-m] <filename>");
            return;
        }

        try
        {
            long lineCount = 0, wordCount = 0, byteCount = 0, charCount = 0;
            bool countLines = false, countWords = false, countBytes = false, countChars = false;

            if (option == "-c") countBytes = true;
            else if (option == "-l") countLines = true;
            else if (option == "-w") countWords = true;
            else if (option == "-m") countChars = true;
            else
            {
                countLines = true;
                countWords = true;
                countBytes = true;
            }

            if (useStdin)
            {
                using (var tempFile = new TempFile())
                {
                    using (var writer = new StreamWriter(tempFile.Path))
                    {
                        string line;
                        while ((line = Console.ReadLine()) != null)
                        {
                            writer.WriteLine(line);
                        }
                    }

                    if (countLines) lineCount = CountLines(tempFile.Path);
                    if (countWords) wordCount = CountWords(tempFile.Path);
                    if (countBytes) byteCount = new FileInfo(tempFile.Path).Length;
                    if (countChars) charCount = CountCharacters(tempFile.Path);
                }

                if (option == "-c") Console.WriteLine($"  {byteCount}");
                else if (option == "-l") Console.WriteLine($"  {lineCount}");
                else if (option == "-w") Console.WriteLine($"  {wordCount}");
                else if (option == "-m") Console.WriteLine($"  {charCount}");
                else Console.WriteLine($"  {lineCount}  {wordCount}  {byteCount}");
            }
            else
            {
                if (!File.Exists(filePath))
                {
                    Console.WriteLine($"Error: File '{filePath}' does not exist.");
                    return;
                }

                if (countLines) lineCount = CountLines(filePath);
                if (countWords) wordCount = CountWords(filePath);
                if (countBytes) byteCount = new FileInfo(filePath).Length;
                if (countChars) charCount = CountCharacters(filePath);

                if (option == "-c") Console.WriteLine($"  {byteCount} {filePath}");
                else if (option == "-l") Console.WriteLine($"  {lineCount} {filePath}");
                else if (option == "-w") Console.WriteLine($"  {wordCount} {filePath}");
                else if (option == "-m") Console.WriteLine($"  {charCount} {filePath}");
                else Console.WriteLine($"  {lineCount}  {wordCount}  {byteCount} {filePath}");
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error: {ex.Message}");
        }
    }

    static long CountLines(string filePath)
    {
        long lineCount = 0;
        using (var reader = new StreamReader(filePath))
        {
            while (reader.ReadLine() != null)
            {
                lineCount++;
            }
        }
        return lineCount;
    }

    static long CountWords(string filePath)
    {
        long wordCount = 0;
        using (var reader = new StreamReader(filePath))
        {
            string line;
            while ((line = reader.ReadLine()) != null)
            {
                if (string.IsNullOrWhiteSpace(line)) continue;
                wordCount += line.Split((char[])null, StringSplitOptions.RemoveEmptyEntries).Length;
            }
        }
        return wordCount;
    }

    static long CountCharacters(string filePath)
    {
        long charCount = 0;
        using (var reader = new StreamReader(filePath, Encoding.UTF8))
        {
            char[] buffer = new char[1024];
            int charsRead;
            while ((charsRead = reader.Read(buffer, 0, buffer.Length)) > 0)
            {
                charCount += charsRead;
            }
        }
        return charCount;
    }
}

class TempFile : IDisposable
{
    public string Path { get; }

    public TempFile()
    {
        Path = System.IO.Path.GetTempFileName();
    }

    public void Dispose()
    {
        if (File.Exists(Path))
        {
            File.Delete(Path);
        }
    }
}