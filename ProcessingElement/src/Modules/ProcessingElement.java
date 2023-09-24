package Modules;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.Arrays;
import java.util.logging.Level;
import java.util.logging.Logger;
/**
 * Represents a processing element hat loads and executes instructions from a program file.
 */
public class ProcessingElement extends Thread {
    private static final Logger logger = Logger.getLogger(ProcessingElement.class.getName());
    private final int coreId; // Unique ID for the PE
    private final String processingElementName; // Name of the processing element
    private int registerRs; // The only register Rs
    private String programFilePath; // The path where the program file is located
    private String[] program; // The abstraction of the ROM memory to store the program
    private int totalInstructions; // The number of instructions loaded in the program array
    private int programCounter; // Program counter to know how many instructions have being executed
    private volatile boolean running; // Flag to know when to stop the Core
    private final GodSharedBuffer buffer; // Shared Buffer to communicate with the Cache Thread
    private ObjectMapper objectMapper;



    private static final String READ_INSTRUCTION_PATTERN = "^(READ (1[0-5]|[0-9]))$";
    private static final String WRITE_INSTRUCTION_PATTERN = "^(WRITE (1[0-5]|[0-9]))$";
    private static final String INC_INSTRUCTION_PATTERN = "^(INC)$";

    /**
     * Initializes a new ProcessingElement instance.
     *
     * @param coreId               The unique ID for the processing element.
     * @param processingElementName The name of the processing element.
     * @param programFilePath      The path to the program file to load instructions from.
     */
    public ProcessingElement(int coreId, String processingElementName, String programFilePath, GodSharedBuffer buffer) {
        this.coreId = coreId;
        this.processingElementName = processingElementName;
        this.programFilePath = programFilePath;
        this.registerRs = 0;
        this.program = new String[]{"", "", "", "", "", ""};
        this.totalInstructions = 0;
        this.programCounter = 0;
        this.running = true;
        this.buffer = buffer;
        this.objectMapper  = new ObjectMapper();
        loadProgramFile();
    }

    /**
     * Reads every line of the program file, validates if they are valid instructions, and stores into the program array.
     */
    private void loadProgramFile() {
        try (FileReader fileReader = new FileReader(programFilePath);
             BufferedReader bufferedReader = new BufferedReader(fileReader)) {

            String instruction;
            int lineNumber = 0;

            while ((instruction = bufferedReader.readLine()) != null && lineNumber < program.length) {
                if (!isValidInstruction(instruction)) {
                    logger.warning("[" + processingElementName + "] Invalid instruction at line " + (lineNumber + 1));
                }

                program[lineNumber] = instruction;
                totalInstructions++;
                lineNumber++;
            }

            logger.info("[" + processingElementName + "] Loaded " + totalInstructions + " instructions successfully.");
            logger.info("[" + processingElementName + "] The instructions loaded are: " + Arrays.toString(program));

        } catch (IOException e) {
            logger.log(Level.SEVERE, "[" + processingElementName + "] An error occurred while reading the file", e);
        }
    }

    /**
     * Validates if an instruction is valid.
     */
    private boolean isValidInstruction(String instruction) {
        return instruction.matches(READ_INSTRUCTION_PATTERN) ||
                instruction.matches(WRITE_INSTRUCTION_PATTERN) ||
                instruction.matches(INC_INSTRUCTION_PATTERN);
    }

    /**
     * Executes the loaded instructions and logs the output.
     */
    private void executeInstruction() {
        for (String instruction : program) {
            if (instruction.matches(READ_INSTRUCTION_PATTERN)) {
                System.out.println("[" + processingElementName + "] READ INSTRUCTION");
            } else if (instruction.matches(WRITE_INSTRUCTION_PATTERN)) {
                System.out.println("[" + processingElementName + "] WRITE INSTRUCTION");
            } else if (instruction.matches(INC_INSTRUCTION_PATTERN)) {
                System.out.println("[" + processingElementName + "] INC INSTRUCTION");
            }
        }
    }

    /**
     * Increments the value stored in the register RS.
     */
    private void incrementRegisterRs() {
        registerRs++;
    }

    /**
     * Returns the value of the Core Identifier.
     */
    public int getCoreId() {
        return coreId;
    }

    /**
     * Executes code in the background when the Thread starts.
     */
    @Override
    public void run() {
        ObjectNode json;
        String jsonString;
        while (running) {
            try {
                // Validate if all instructions have been executed already
                if (programCounter == totalInstructions){
                    running = false;
                    logger.info("[" + processingElementName + "] has executed all the instructions.");
                    break;
                }
                // Load the instruction from the program array
                String instruction = program[programCounter];
                //logger.info("[" + processingElementName + "] is executing the instruction: " + instruction);
                String instructionType = "";
                if (instruction.matches(READ_INSTRUCTION_PATTERN)) {
                    instructionType = "READ";

                } else if (instruction.matches(WRITE_INSTRUCTION_PATTERN)) {
                    instructionType = "WRITE";

                } else if (instruction.matches(INC_INSTRUCTION_PATTERN)) {
                    programCounter += 1;
                    continue;
                }

                // Create a JSON object to package the message
                json = objectMapper.createObjectNode();
                json.put("instructionType", instructionType);
                json.put("data", 30);

                // Convert JSON object to String
                jsonString = objectMapper.writeValueAsString(json);
                logger.info("[" + processingElementName + "] is ready to send a request to the Cache Controller.");

                // Send request to the cache controller **************************************************************
                buffer.produce(jsonString, true);
                logger.info("[" + processingElementName + "] sent a request to the cache controller.");

                String response = buffer.consume(true);
                logger.info("[" + processingElementName + "] Received the response: " + response);

                // Increment the program counter
                programCounter += 1;

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break; // Exit the loop on InterruptedException
            } catch (JsonProcessingException e) {
                throw new RuntimeException(e);
            }
        }
    }
}
