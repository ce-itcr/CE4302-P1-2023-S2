import Modules.CacheController;
import Modules.ProcessingElement;
import Modules.GodSharedBuffer;

public class Main {
    public static void main(String[] args) {
        // Shared Buffer to communicate the PE and the CacheController
        GodSharedBuffer buffer = new GodSharedBuffer(1, "PE&Cache");

        // Instantiate a PE
        ProcessingElement PE0 = new ProcessingElement(1, "PE0", "src/programFiles/program1.txt", buffer);
        CacheController CC = new CacheController(buffer);

        // Start the Threads
        PE0.start();
        CC.start();

    }
}