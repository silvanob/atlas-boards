/*
 * Copyright 2015 The gRPC Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package me.biemans.atlasboards;

import com.google.common.annotations.VisibleForTesting;
import com.google.protobuf.Message;
import io.grpc.Channel;
import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.StreamObserver;
import me.biemans.atlasboards.routeguide.*;

import java.io.IOException;
import java.util.Random;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * Sample client code that makes gRPC calls to the server.
 */
public class RouteGuideClient {
    private static final Logger logger = Logger.getLogger(RouteGuideClient.class.getName());

    private final RouteGuideGrpc.RouteGuideBlockingStub blockingStub;
    private final RouteGuideGrpc.RouteGuideStub asyncStub;

    private Random random = new Random();
    private TestHelper testHelper;

    /** Construct client for accessing RouteGuide server using the existing channel. */
    public RouteGuideClient(Channel channel) {
        blockingStub = RouteGuideGrpc.newBlockingStub(channel);
        asyncStub = RouteGuideGrpc.newStub(channel);
    }

    /**
     * Blocking unary call example.  Calls getFeature and prints the response.
     */
    public void getFeature(String lat, String lon) {
        info("*** GetFeature: lat={0} lon={1}", lat, lon);

        Ticket ticket = Ticket.newBuilder().setTitle("Hello").setContent("Bye!!").build();

        Ticket ticketReceived;
        try {
            ticketReceived = blockingStub.createTicket(ticket);
            if (testHelper != null) {
                testHelper.onMessage(ticketReceived);
            }
        } catch (StatusRuntimeException e) {
            warning("RPC failed: {0}", e.getStatus());
            if (testHelper != null) {
                testHelper.onRpcError(e);
            }
            return;
        }
        //if (RouteGuideUtil.exists(feature)) {
        //    info("Found feature called \"{0}\" at {1}, {2}",
        //            feature.getName(),
        //            RouteGuideUtil.getLatitude(feature.getLocation()),
        //            RouteGuideUtil.getLongitude(feature.getLocation()));
        //} else {
        //    info("Found no feature at {0}, {1}",
        //            RouteGuideUtil.getLatitude(feature.getLocation()),
        //            RouteGuideUtil.getLongitude(feature.getLocation()));
        //}
    }

    /**
     * Blocking server-streaming example. Calls listFeatures with a rectangle of interest. Prints each
     * response feature as it arrives.
     */

    /**
     * Async client-streaming example. Sends {@code numPoints} randomly chosen points from {@code
     * features} with a variable delay in between. Prints the statistics when they are sent from the
     * server.
     */

    /** Issues several different requests and then exits. */
    public static void main(String[] args) throws InterruptedException {
        String target = "localhost:50051";
        if (args.length > 0) {
            if ("--help".equals(args[0])) {
                System.err.println("Usage: [target]");
                System.err.println("");
                System.err.println("  target  The server to connect to. Defaults to " + target);
                System.exit(1);
            }
            target = args[0];
        }

        //List<Feature> features;
        //try {
        //    features = RouteGuideUtil.parseFeatures(RouteGuideUtil.getDefaultFeaturesFile());
        //} catch (IOException ex) {
        //    ex.printStackTrace();
        //    return;
        //}

        ManagedChannel channel = Grpc.newChannelBuilder(target, InsecureChannelCredentials.create())
                .build();
        try {
            RouteGuideClient client = new RouteGuideClient(channel);
            // Looking for a valid feature
            client.getFeature("general", "world");

//            // Feature missing.
//            client.getFeature(0, 0);
//
//            // Looking for features between 40, -75 and 42, -73.
//            client.listFeatures(400000000, -750000000, 420000000, -730000000);
//
//            // Record a few randomly selected points from the features file.
//            client.recordRoute(features, 10);

            // Send and receive some notes.
//            CountDownLatch finishLatch = client.routeChat();

//            if (!finishLatch.await(1, TimeUnit.MINUTES)) {
//                client.warning("routeChat can not finish within 1 minutes");
//            }
        } finally {
            channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
        }
    }

    private void info(String msg, Object... params) {
        logger.log(Level.INFO, msg, params);
    }

    private void warning(String msg, Object... params) {
        logger.log(Level.WARNING, msg, params);
    }


    /**
     * Only used for unit test, as we do not want to introduce randomness in unit test.
     */
    @VisibleForTesting
    void setRandom(Random random) {
        this.random = random;
    }

    /**
     * Only used for helping unit test.
     */
    @VisibleForTesting
    interface TestHelper {
        /**
         * Used for verify/inspect message received from server.
         */
        void onMessage(Message message);

        /**
         * Used for verify/inspect error received from server.
         */
        void onRpcError(Throwable exception);
    }

    @VisibleForTesting
    void setTestHelper(TestHelper testHelper) {
        this.testHelper = testHelper;
    }
}
