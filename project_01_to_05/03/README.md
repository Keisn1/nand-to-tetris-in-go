
# Table of Contents

1.  [Memory](#orga98a331)
    1.  [Sequential Logic](#org8abf6a0)
        1.  [Time steps](#orgf74431b)
        2.  [Technical implementation](#orge35d319)
        3.  [Latch / Data Flip Flop](#org4e2e2b2)
    2.  [Project](#org341bb37)
        1.  [Bit](#org2a92a28)
        2.  [PC (Program Counter)](#orgee3be4d)



<a id="orga98a331"></a>

# Memory


<a id="org8abf6a0"></a>

## Sequential Logic

Up until this point we were concerned with Logic Gates and Chips that are time independent, meaning that they are only dependent on the current inputs. This is called **Combination Logic**.

We now like to introduce **Sequential Logic** where the output depends on the previous inputs as well.


<a id="orgf74431b"></a>

### Time steps

Time will be represented as a discrete sequence of time ($$t = 0, 1, 2, 3, 4 ...$$).


<a id="orge35d319"></a>

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


<a id="org4e2e2b2"></a>

### Latch / Data Flip Flop

The most elemetary sequential gate that we can imagine is called a *Latch* or *Data Flip Flop*. The *Latch* outputs the input of the previous Time step.

![img](imgs/dff.png)

In the course, this DFF is given to us as a pre-implemented Chip that we can use. From this we are building a Bit and then all the other Chips ([Project](#org341bb37)).


<a id="org341bb37"></a>

## Project

In the project we are building the memory and a program counter. In detail, we are implementing the following Chips using our **Hardware Description Language**.

1.  `Bit, Register, RAM8, RAM64` and `PC` (Program Counter) in directory [a](https://github.com/Keisn1/nand-to-tetris-in-go/tree/main/project_01_to_05/03/a)
2.  `RAM512, RAM4k, RAM16K` in directory [b](https://github.com/Keisn1/nand-to-tetris-in-go/tree/main/project_01_to_05/03/b)

The central Chip is th `Bit`. All other Memory Chips are built from the `Bit` (`Register`, `RAM8` &#x2026; ).


<a id="org2a92a28"></a>

### Bit

To build a `Bit`, one needs to realize the behavior of &ldquo;Loading&rdquo; and &ldquo;Storing&rdquo;. This is done using a **DFF** in combination with a `MUX` the takes as **sel** a **load** Flag.

![img](imgs/bit.png)


<a id="orgee3be4d"></a>

### PC (Program Counter)

The `PC` will be a `Register` that takes additional Flags, so that we can **Reset** (to 0), **Increment** (next instruction in order) or **Load** (Jump to another to another given instruction = input value) the next instruction to be executed.

![img](imgs/pc.png)

