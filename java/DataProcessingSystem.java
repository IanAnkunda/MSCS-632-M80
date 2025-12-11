import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.util.LinkedList;
import java.util.List;
import java.util.Queue;

// Simple shared queue that multiple threads can use safely
class SharedQueue {
    private final Queue<String> queue = new LinkedList<>();

    // Add a task and wake up any waiting workers
    public synchronized void addTask(String task) {
        queue.add(task);
        notifyAll();
    }

    // Workers call this to get a task; wait if none available
    public synchronized String getTask() throws InterruptedException {
        while (queue.isEmpty()) {
            wait();
        }
        return queue.poll();
    }
}

// Worker thread that keeps pulling tasks, "processing" them, and saving results
class Worker extends Thread {
    private final SharedQueue queue;
    private final List<String> results;

    public Worker(SharedQueue queue, List<String> results) {
        this.queue = queue;
        this.results = results;
    }

    @Override
    public void run() {
        try {
            while (true) {
                String task = queue.getTask();

                // simulate some work
                Thread.sleep(200);

                String output = getName() + " processed " + task;

                // save result safely
                synchronized (results) {
                    results.add(output);
                }
            }
        } catch (InterruptedException e) {
            // thread stops normally when interrupted
        }
    }
}

public class DataProcessingSystem {
    public static void main(String[] args) {
        SharedQueue sharedQueue = new SharedQueue();
        List<String> results = new LinkedList<>();

        int workerCount = 4;
        Worker[] workers = new Worker[workerCount];

        // start the workers
        for (int i = 0; i < workerCount; i++) {
            workers[i] = new Worker(sharedQueue, results);
            workers[i].start();
        }

        // add some tasks
        for (int i = 1; i <= 10; i++) {
            sharedQueue.addTask("Task-" + i);
        }

        // give workers time to finish
        try {
            Thread.sleep(2000);
        } catch (InterruptedException ignored) {}

        // stop workers
        for (Worker worker : workers) {
            worker.interrupt();
        }

        // write results to file
        try (BufferedWriter writer = new BufferedWriter(new FileWriter("results_java.txt"))) {
            for (String line : results) {
                writer.write(line);
                writer.newLine();
            }
        } catch (IOException e) {
            System.out.println("File error: " + e.getMessage());
        }
    }
}
