package Modules;

import java.util.LinkedList;
import java.util.Queue;

public class GodSharedBuffer {
    private Queue<String> buffer = new LinkedList<>();
    private int capacity;
    private final String bufferName;
    private boolean isThreadATurn = true;

    public GodSharedBuffer(int capacity, String bufferName) {
        this.capacity = capacity;
        this.bufferName = bufferName;
    }

    public synchronized void produce(String item, boolean isThreadA) throws InterruptedException {
        while ((isThreadA && !isThreadATurn) || (!isThreadA && isThreadATurn) || buffer.size() == capacity) {
            wait();
        }
        buffer.add(item);
        System.out.println("[" + bufferName + "]: A new message was produced in the buffer...");
        isThreadATurn = !isThreadA;
        notify(); // Notify the waiting thread
    }

    public synchronized String consume(boolean isThreadA) throws InterruptedException {
        while ((isThreadA && isThreadATurn) || (!isThreadA && !isThreadATurn) || buffer.isEmpty()) {
            wait();
        }
        String item = buffer.poll();
        System.out.println("[" + bufferName + "]: Item has been consumed...");
        isThreadATurn = !isThreadA;
        notify(); // Notify the waiting thread
        return item;
    }
}




