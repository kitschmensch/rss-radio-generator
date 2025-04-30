document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("rssForm");
  const addStationButton = document.getElementById("addStation");
  const generatedURLInput = document.getElementById("generatedURL");
  const copyURLButton = document.getElementById("copyURL");
  const stationsContainer = document.getElementById("stations");

  // Helper function to update the URL query parameters using base64 encoding
  function updateQueryParams() {
    // Get raw values
    const title = form.title.value;
    const description = form.description.value;

    const stations = Array.from(stationsContainer.querySelectorAll(".station"))
      .map(station => {
        const stationTitle = station.querySelector('input[name="stationTitle"]').value;
        const stationURL = station.querySelector('input[name="stationURL"]').value;
        const stationDescription = station.querySelector('input[name="stationDescription"]').value;
        
        // Only create station entries for valid inputs
        if (!stationTitle && !stationURL && !stationDescription) return null;
        
        return {
          title: stationTitle,
          url: stationURL,
          description: stationDescription
        };
      })
      .filter(station => station !== null);

    // Create a single object with all form data
    const formData = {
      title,
      description,
      stations
    };
    
    // Convert to JSON and encode as base64
    const jsonString = JSON.stringify(formData);
    let base64Data = btoa(unescape(encodeURIComponent(jsonString))); // Handles UTF-8 characters
    
    // Make base64 URL safe by replacing + with - and / with _
    base64Data = base64Data.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
    
    // Create URL with single parameter
    const queryParams = new URLSearchParams();
    queryParams.set("data", base64Data);
    
    // Update the browser URL without reloading the page
    const pageURL = `${window.location.pathname}?${queryParams.toString()}`;
    window.history.replaceState(null, "", pageURL);
    
    // Generate the RSS URL
    const rssURL = `${window.location.origin}/rss?${queryParams.toString()}`;
    
    // Update the "Generated URL" input box
    generatedURLInput.value = rssURL;
  }

  // Populate form fields from base64 encoded data
  function populateFormFromQueryParams() {
    const params = new URLSearchParams(window.location.search);
    const base64Data = params.get("data");
    
    if (!base64Data) return;
    
    try {
      // Handle the URL-safe base64 by replacing URL-safe chars and adding padding if needed
      let sanitizedBase64 = base64Data.replace(/-/g, '+').replace(/_/g, '/');
      while (sanitizedBase64.length % 4) {
        sanitizedBase64 += '=';
      }
      
      // Decode the base64 data
      const jsonString = decodeURIComponent(escape(atob(sanitizedBase64)));
      const formData = JSON.parse(jsonString);
      
      console.log("Decoded form data:", formData); // Debug log
      
      // Populate form fields
      if (formData.title) form.title.value = formData.title;
      if (formData.description) form.description.value = formData.description;
      
      // Add stations
      if (formData.stations && Array.isArray(formData.stations) && formData.stations.length > 0) {
        stationsContainer.innerHTML = ""; // Clear existing stations
        
        formData.stations.forEach(station => {
          if (!station) return; // Skip null/undefined stations
          
          const stationDiv = document.createElement("div");
          stationDiv.classList.add("station");
          
          // Sanitize values to prevent XSS
          const safeTitle = (station.title || '').replace(/</g, '&lt;').replace(/"/g, '&quot;');
          const safeUrl = (station.url || '').replace(/</g, '&lt;').replace(/"/g, '&quot;');
          const safeDescription = (station.description || '').replace(/</g, '&lt;').replace(/"/g, '&quot;');
          
          stationDiv.innerHTML = `
            <input type="text" name="stationTitle" placeholder="Station Title" value="${safeTitle}" required>
            <input type="url" name="stationURL" placeholder="Station URL" value="${safeUrl}" required>
            <input type="text" name="stationDescription" placeholder="Station Description" value="${safeDescription}" required>
            <button type="button" class="removeStation">Remove</button>
          `;
          stationsContainer.appendChild(stationDiv);
          
          // Add event listener to remove button
          stationDiv.querySelector(".removeStation").addEventListener("click", () => {
            stationDiv.remove();
            updateQueryParams();
          });
        });
      }
    } catch (e) {
      console.error("Error decoding form data:", e);
      // Optionally display error to user
    }
  }

  // Add a new station input group
  addStationButton.addEventListener("click", () => {
    const stationDiv = document.createElement("div");
    stationDiv.classList.add("station");
    stationDiv.innerHTML = `
      <input type="text" name="stationTitle" placeholder="Station Title" required>
      <input type="url" name="stationURL" placeholder="Station URL" required>
      <input type="text" name="stationDescription" placeholder="Station Description" required>
      <button type="button" class="removeStation">Remove</button>
    `;
    stationsContainer.appendChild(stationDiv);

    // Add event listener to remove button
    stationDiv.querySelector(".removeStation").addEventListener("click", () => {
      stationDiv.remove();
      updateQueryParams();
    });

    updateQueryParams();
  });

  // Update query parameters and generated URL on form input
  form.addEventListener("input", updateQueryParams);

  // Copy the generated URL to the clipboard
  copyURLButton.addEventListener("click", () => {
    generatedURLInput.select();
    document.execCommand("copy");
    alert("URL copied to clipboard!");
  });

  // Populate the form with query parameters on page load
  populateFormFromQueryParams();

  // Generate the initial URL
  updateQueryParams();
});