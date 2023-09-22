import Modules.ProcessingElement;

public class Main {
    public static void main(String[] args) {
        // Instantiate a PE
        ProcessingElement PE0 = new ProcessingElement(1, "PE0", "src/programFiles/program1.txt");
        PE0.start();
    }
}