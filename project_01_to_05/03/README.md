
# Table of Contents

1.  [Memory](#org56705f9)
    1.  [Sequential Logic](#orged5821b)
        1.  [Time steps](#org353453a)
        2.  [Technical implementation](#org8fdccf5)
        3.  [Latch / Data Flip Flop](#org09b17be)
    2.  [Project](#org7ebd31f)



<a id="org56705f9"></a>

# Memory


<a id="orged5821b"></a>

## Sequential Logic

Up until this point we were concerned with Logic Gates and Chips that are time independent, meaning that they are only dependent on the current inputs. This is called **Combination Logic**.

We now like to introduce **Sequential Logic** where the output depends on the previous inputs as well.


<a id="org353453a"></a>

### Time steps

Time will be represented as a discrete sequence of time ($$t = 0, 1, 2, 3, 4 ...$$).


<a id="org8fdccf5"></a>

### Technical implementation

To build a clock with discrete time units, we are going to use an oscillator, which is oscillating between on and off at a given rate ($$ freq=[H] $$).

In Each time unit, any input or output value of a gate is continious.

![img](imgs/clock_ideal.png)

1.  Physical problems

    The change in current and voltage is continuous as well. Therefore there will be delays in the signal of the oscillator and the signal will rather look like the following.
    
    ![img](imgs/clock_real.png)
    
    These delays need to be taken into account and a common design decision is to track state at the end of a cycle.
    
    ![img](imgs/clock_design.png)
    
    The resulting effect is that our sequential logic gates are only reacting to a given input at each end of a cycle, whereas combinational gates are going to &ldquo;react&rdquo; immediately.


<a id="org09b17be"></a>

### Latch / Data Flip Flop

The most elemetary sequential gate that we can imagine is called a *Latch* or *Data Flip Flop*. The *Latch* outputs the input of the previous Time step.

![img](imgs/dff.png)

In the course, this DFF is given to us as a pre-implemented Chip that we can use. From this we are building a Bit and then all the other Chips ([Project](#org7ebd31f)).


<a id="org7ebd31f"></a>

## Project

In the project we are building the memory and a program counter. In detail, we are implementing the following Chips using our **Hardware Description Language**.

1.  `Bit, Register, RAM8, RAM64` and `PC` (Program Counter)
2.  `RAM512, RAM4k, RAM16K`

