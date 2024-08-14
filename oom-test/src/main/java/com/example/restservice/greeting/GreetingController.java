package com.example.restservice.greeting;

import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class GreetingController {

    private static final String     template = "Hello, %s!";
    private final        AtomicLong counter  = new AtomicLong();

    @GetMapping("/greeting")
    public Greeting greeting(@RequestParam(value = "name", defaultValue = "World") String name) {
        return new Greeting(counter.incrementAndGet(), String.format(template, name));
    }

    @GetMapping("/oom")
    public String oom() {
        List<ByteBuffer> list = new ArrayList<>();
        while (true) {
            list.add(ByteBuffer.allocate(10_000_000));
        }
    }

}
