import contextlib

# Global state
_global_state = {}


@contextlib.contextmanager
def update_global_context(**kwargs):
    """Create a global context and use it across the project."""
    global _global_state  # Declare that we're using the global state

    # Update _global_state with new key-value pairs
    _global_state.update(kwargs)


@contextlib.contextmanager
def delete_global_context(key):
    # Clean up the global state
    del _global_state[key]
