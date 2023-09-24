package Modules;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;

public class CacheController extends Thread {
    private final GodSharedBuffer buffer;
    private volatile boolean running;

    public CacheController(GodSharedBuffer buffer) {
        this.buffer = buffer;
        this.running = true;
    }

    public boolean getRunning() {
        return running;
    }

    public void setRunning(boolean newRunning) {
        running = newRunning;
    }

    @Override
    public void run() {
        // Parse the JSON String back to a JSON object
        ObjectMapper objectMapper = new ObjectMapper();
        while (running) {
            try {
                // Wait for a request from ProcessingElement
                String request = buffer.consume(false);
                ObjectNode parsedJson = objectMapper.readValue(request, ObjectNode.class);

                // Access values from the parsed JSON object
                String instructionType = parsedJson.get("instructionType").asText();
                int data = parsedJson.get("data").asInt();

                // Display the data received
                //System.out.println("Instruction: " + instructionType + " and data: " + data);
                // Simulate processing time (adjust as needed)
                Thread.sleep(10000);

                buffer.produce("Response from Cache", false);

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break; // Exit the loop on InterruptedException
            } catch (JsonProcessingException e) {
                throw new RuntimeException(e);
            }
        }
    }
}

