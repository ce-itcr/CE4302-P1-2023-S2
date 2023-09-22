package Modules;
import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProcessingElement extends Thread {
    private static final Logger logger = Logger.getLogger(ProcessingElement.class.getName());
    private final int core_id; // Unique ID for the PE
    //private Cache cache; // Reference to the private cache
    private int registerRs; // The only register Rs
    private String[] program; // The abstraction of the ROM memory to store the program
    private int programCounter; // The PC to control the instruction count

    // Constructor
    public ProcessingElement(int core_id) {
        this.core_id = core_id;
        //this.cache = cache;
        this.registerRs = 0; // Initialize Rs to 0
        this.program = new String[]{"", "", "", ""}; // Initialize the program array to empty
        this.programCounter = 0; // Initialize the PC to 0
    }

    // Method to read a text file containing the program to be executed
    public void ReadProgramFile(String programPath) {
        System.out.println("Reading program at: " + programPath);
        try {
            // Create a FileReader and BufferedReader to read the file
            FileReader fileReader = new FileReader(programPath);
            BufferedReader bufferedReader = new BufferedReader(fileReader);

            String line;
            // Declare a counter to know the line number
            int lineNumber = 0;

            // Read and print each line from the file
            while ((line = bufferedReader.readLine()) != null) {
                System.out.println(lineNumber + ": " + line);
                lineNumber += 1;
            }

            // Close the BufferedReader and FileReader when done
            bufferedReader.close();
            fileReader.close();
        } catch (IOException e) {
            // Handle any potential exceptions, e.g., file not found or unable to read
            logger.log(Level.SEVERE, "An error occurred while reading the file", e);
        }
    }

    // Method to load instructions from a program file into the program attribute
    private boolean loadProgramFile(String programFilePath){
        try {
            // Create a FileReader and BufferedReader to read the file
            FileReader fileReader = new FileReader(programFilePath);
            BufferedReader bufferedReader = new BufferedReader(fileReader);

            // String to load every line
            String instruction;
            // Counter to know the line number
            int lineNumber = 0;

            // Read and print each line from the file
            while ((instruction = bufferedReader.readLine()) != null) {
                // Check if the instructions count has exceeded the maximum number
                if (lineNumber > 3){
                    System.out.println("Error: The program loaded at " + programFilePath + " exceeds the maximum number" +
                            " of instructions");
                    return false;
                }

                // Check if the instruction is valid
                if (!isValidInstruction(instruction)){
                    System.out.println("The instruction at line " + lineNumber + " is invalid");
                    return false;
                }

                // If there is no problem with the instruction, then, store it in the program attribute
                program[lineNumber] = instruction;

                // Increment the line number
                lineNumber += 1;
                }

            // Close the BufferedReader and FileReader when done
            bufferedReader.close();
            fileReader.close();
            return true;
        } catch (IOException e) {
            // Handle any potential exceptions, e.g., file not found or unable to read
            logger.log(Level.SEVERE, "An error occurred while reading the file", e);
            return false;
        }
    }

    // Method to validate if an instruction is valid
    private boolean isValidInstruction(String instruction) {
        // Define a regular expression pattern for the allowed instructions
        String pattern = "^(READ \\d+|Write (1[0-5]|[0-9])|INC)$";

        // Compile the pattern
        Pattern regex = Pattern.compile(pattern);

        // Match the input string against the pattern
        Matcher matcher = regex.matcher(instruction);

        // Return true if there is a match, indicating a valid instruction
        return matcher.matches();
    }

    // Method to execute an instruction
    public void executeInstruction(int instructionType, int address, int data) {
        switch (instructionType) {
            case 1: // Read an address from memory/cache and store it in Rs
                readMemory(address);
                break;
            case 2: // Write the value stored in Rs in memory/cache
                writeMemory(address, registerRs);
                break;
            case 3: // Increment the value of Rs
                incrementRegisterRs();
                break;
            default:
                System.out.println("Invalid instruction type.");
        }
    }

    // Method to read an address from memory/cache and store it in Rs
    private void readMemory(int address) {
        //int data = cache.read(address); // Read from the cache
        registerRs = 100; // Store in Rs
    }

    // Method to write the value stored in Rs in memory/cache
    private void writeMemory(int address, int data) {
        //cache.write(address, data); // Write to the cache
        System.out.println("Writing the value: " + data + " at the address: " + address);
    }

    // Method to increment the value of Rs
    private void incrementRegisterRs() {
        registerRs++;
    }

    // Getter for ID
    public int getCore_id() {
        return core_id;
    }

    // Getter for Rs register value
    public int getRegisterRs() {
        return registerRs;
    }

    // Setter for Rs register value
    public void setRegisterRs(int value) {
        this.registerRs = value;
    }

    // Getter for the cache
//    public Cache getCache() {
//        return cache;
//    }

    @Override
    public void run() {
        // You can implement the logic for the PE's operation here,
        // including executing a sequence of instructions.
    }
}

