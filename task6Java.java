import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.Scanner;
import java.util.HashMap;
import java.util.Map;


class Tunnel {
    private final String tunnelName;
    private final BlockingQueue<String> trainQueue;

    public Tunnel(String tunnelName) {
        this.tunnelName = tunnelName;
        this.trainQueue = new ArrayBlockingQueue<>(10);
    }

    public void enterTunnel(String trainName) throws InterruptedException {
        trainQueue.put(trainName);
        System.out.println(trainName + " очікує на в'їзд у " + tunnelName);

        String oppositeTunnelName = tunnelName.equals("Тунель 1") ? "Тунель 2" : "Тунель 1";
        Tunnel oppositeTunnel = Main.getTunnel(oppositeTunnelName);

        if (!oppositeTunnel.trainQueue.isEmpty()) {
            String oppositeTrain = oppositeTunnel.trainQueue.peek();
            System.out.println(trainName + " перевіряє чи є " + oppositeTrain + " в протилежному тунелі.");

            if (Main.hasExceededWaitTime(oppositeTrain)) {
                System.out.println(trainName + " змінює рух до " + tunnelName);
                oppositeTunnel.trainQueue.take();
                trainQueue.put(oppositeTrain);
            }
        }

        System.out.println(trainName + " проїжджає " + tunnelName + " тунель.");
        Thread.sleep(10000);
        System.out.println(trainName + " проїхав " + tunnelName + " тунель.");
        trainQueue.take();
        Main.incrementCounter(tunnelName);
    }
}

class Train extends Thread {
    private final int trainNumber;
    private final Tunnel tunnel;

    public Train(int trainNumber, Tunnel tunnel) {
        this.trainNumber = trainNumber;
        this.tunnel = tunnel;
    }

    @Override
    public void run() {
        try {
            tunnel.enterTunnel("Потяг " + trainNumber);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }
}

public class Main {
    private static final Tunnel tunnel1 = new Tunnel("Тунель 1");
    private static final Tunnel tunnel2 = new Tunnel("Тунель 2");
    private static int trainsPassedTunnel1 = 0;
    private static int trainsPassedTunnel2 = 0;
    private static final Map<String, Long> trainArrivalTimes = new HashMap<>();

    public static Tunnel getTunnel(String tunnelName) {
        return tunnelName.equals("Тунель 1") ? tunnel1 : tunnel2;
    }

    public static synchronized void incrementCounter(String tunnelName) {
        if (tunnelName.equals("Тунель 1")) {
            trainsPassedTunnel1++;
        } else if (tunnelName.equals("Тунель 2")) {
            trainsPassedTunnel2++;
        }

        int totalTrains = trainsPassedTunnel1 + trainsPassedTunnel2;
        if (totalTrains == 20) {
            System.out.println("Колії вільні");
        }
    }

    public static boolean hasExceededWaitTime(String trainName) {
        if (!trainArrivalTimes.containsKey(trainName)) {
            trainArrivalTimes.put(trainName, System.currentTimeMillis());
            return false;
        } else {
            long arrivalTime = trainArrivalTimes.get(trainName);
            long currentTime = System.currentTimeMillis();
            long waitingTime = currentTime - arrivalTime;
            return waitingTime > 60000; // Перевірка, чи час очікування перевищує 60 секунд
        }
    }

    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        System.out.println("Яку кількість потягів очікуємо?");
        int totalTrains = scanner.nextInt();

        for (int i = 1; i <= totalTrains; i++) {
            Tunnel selectedTunnel = i % 2 == 0 ? tunnel1 : tunnel2;
            Train train = new Train(i, selectedTunnel);
            train.start();
        }
    }
}
