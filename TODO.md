# Tui

## Todo

- Quicksort:
    - [ ] add highlight where current i (number of values less than pivot) is
    - [ ] fix locked highlighting

- Heapsort:
    - [ ] add locked
    - [ ] coloring for paths
    - [ ] add !interesting (Lightgray)
    - [ ] show current state of array in array style
    - [ ] center description
    - [ ] Coloring and Updating needs more work, sometimes there remain relics
      from previous state

- Search:
    - [ ] Always show what is the target number (Box?)
    - [ ] Don't show sorted array as last step -> show if target was found
    (highlight) or was not found
    - [ ] have different starting and ending columnGraphData for search and sort

- [ ] make statusbar smarter when After() 3s and meanwhile set() then set() gets
        removed after tick
- [ ] Layout adjustments based on window size
- [ ] Option to sort ascending or descending
- [ ] Option between column or array display
- [ ] Clean up main file
- [ ] Right now every widget returns the string instructions for buf, make it
  that either buf is injected into these methods or there is a channel where
  buf listens on, ...

- Box
    - [ ] handle if title is longer than content
    - [ ] handle const size, so that if the content shrinks the box size stays
      the same

<!-- FIX:
- [ ] move: name, graphState and controlCh into columnGraph ? -> How to do it?
-->

- [ ] improve coloring methods for adding Bold, BG, FG colors
