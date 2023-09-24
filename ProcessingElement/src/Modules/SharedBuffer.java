package Modules;

import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.Semaphore;

public class SharedBuffer {
    private Queue<String> buffer = new LinkedList<>();
    private Semaphore requestSemaphore = new Semaphore(0);
    private Semaphore responseSemaphore = new Semaphore(0);
    private int capacity;
    private boolean responseReady = false;

    public SharedBuffer(int capacity) {
        this.capacity = capacity;
    }

    public synchronized void sendRequest(String request) throws InterruptedException {
        while (buffer.size() == capacity) {
            System.out.println("Buffer is full, waiting for space...");
            wait(); // Wait while the buffer is full
        }
        buffer.add(request);
        System.out.println("Sent request: " + request);
        requestSemaphore.release();
    }

    public synchronized String processRequest() throws InterruptedException {
        while (buffer.isEmpty()) {
            System.out.println("Buffer is empty, waiting for a request...");
            wait(); // Wait while the buffer is empty
        }
        requestSemaphore.acquire();
        String request = buffer.poll();
        //System.out.println("Processing request: " + request);
        responseReady = true;
        return request;
    }

    public synchronized void sendResponse(String response) throws InterruptedException {
        while (!responseReady) {
            System.out.println("Waiting for response...");
            wait(); // Wait until a response is expected
        }
        buffer.add(response);
        System.out.println("Sent response: " + response);
        responseSemaphore.release();
        responseReady = false;
        notify(); // Notify that a response is sent
    }

    public synchronized String waitForResponse() throws InterruptedException {
        while (!responseReady) {
            System.out.println("Waiting for response...");
            wait(); // Wait until a response is ready
        }
        responseSemaphore.acquire();
        String response = buffer.poll();
        System.out.println("Received response: " + response);
        notify(); // Notify that a response is received
        responseReady = false;
        return response;
    }
}
