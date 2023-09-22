import Modules.ProcessingElement;

public class Main {
    public static void main(String[] args) {
        // Instantiate a PE
        ProcessingElement PE0 = new ProcessingElement(1);
        // Make the PE read the program file
        PE0.ReadProgramFile("src/programFiles/program1.txt");
    }
}